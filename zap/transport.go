// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrClosed         = errors.New("zap: connection closed")
	ErrTimeout        = errors.New("zap: request timeout")
	ErrResponseFailed = errors.New("zap: response failed")
)

// Config contains transport configuration
type Config struct {
	// ReadTimeout is the timeout for reading a message
	ReadTimeout time.Duration
	// WriteTimeout is the timeout for writing a message
	WriteTimeout time.Duration
	// MaxConcurrent is the maximum number of concurrent requests
	MaxConcurrent int
	// BufferSize is the read/write buffer size
	BufferSize int
}

// DefaultConfig returns a Config with reasonable defaults
func DefaultConfig() *Config {
	return &Config{
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  10 * time.Second,
		MaxConcurrent: 1000,
		BufferSize:    64 * 1024,
	}
}

// Conn is a ZAP connection that multiplexes requests
type Conn struct {
	conn   net.Conn
	config *Config

	reader *bufio.Reader
	writer *bufio.Writer

	mu       sync.Mutex
	requests map[uint32]chan *response
	nextID   uint32

	closed atomic.Bool
	done   chan struct{}
}

type response struct {
	msgType MessageType
	payload []byte
	err     error
}

// Dial connects to a ZAP server
func Dial(ctx context.Context, addr string, config *Config) (*Conn, error) {
	if config == nil {
		config = DefaultConfig()
	}

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("zap: dial failed: %w", err)
	}

	return newConn(conn, config), nil
}

// NewConn wraps an existing net.Conn as a ZAP connection
func NewConn(conn net.Conn, config *Config) *Conn {
	if config == nil {
		config = DefaultConfig()
	}
	return newConn(conn, config)
}

func newConn(conn net.Conn, config *Config) *Conn {
	c := &Conn{
		conn:     conn,
		config:   config,
		reader:   bufio.NewReaderSize(conn, config.BufferSize),
		writer:   bufio.NewWriterSize(conn, config.BufferSize),
		requests: make(map[uint32]chan *response),
		done:     make(chan struct{}),
	}
	go c.readLoop()
	return c
}

// Close closes the connection
func (c *Conn) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	close(c.done)

	c.mu.Lock()
	for _, ch := range c.requests {
		ch <- &response{err: ErrClosed}
	}
	c.requests = nil
	c.mu.Unlock()

	return c.conn.Close()
}

// Call sends a request and waits for a response
func (c *Conn) Call(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error) {
	if c.closed.Load() {
		return 0, nil, ErrClosed
	}

	// Allocate request ID and response channel
	id := atomic.AddUint32(&c.nextID, 1)
	respCh := make(chan *response, 1)

	c.mu.Lock()
	if c.requests == nil {
		c.mu.Unlock()
		return 0, nil, ErrClosed
	}
	c.requests[id] = respCh
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.requests, id)
		c.mu.Unlock()
	}()

	// Write request with request ID prefix
	reqBuf := GetBuffer()
	defer PutBuffer(reqBuf)

	reqBuf.WriteUint32(id)
	reqBuf.Grow(reqBuf.Len() + len(payload))
	copy(reqBuf.Data[reqBuf.offset:], payload)
	reqBuf.offset += len(payload)

	if err := c.write(msgType, reqBuf.Bytes()); err != nil {
		return 0, nil, err
	}

	// Wait for response
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	case <-c.done:
		return 0, nil, ErrClosed
	case resp := <-respCh:
		if resp.err != nil {
			return 0, nil, resp.err
		}
		return resp.msgType, resp.payload, nil
	}
}

// Send sends a one-way message (no response expected)
func (c *Conn) Send(msgType MessageType, payload []byte) error {
	if c.closed.Load() {
		return ErrClosed
	}
	return c.write(msgType, payload)
}

func (c *Conn) write(msgType MessageType, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.config.WriteTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
	}

	if err := WriteMessage(c.writer, msgType, payload); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Conn) readLoop() {
	for {
		if c.config.ReadTimeout > 0 {
			_ = c.conn.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
		}

		msgType, payload, err := ReadMessage(c.reader)
		if err != nil {
			if !c.closed.Load() {
				c.Close()
			}
			return
		}

		fmt.Fprintf(os.Stderr, "[ZAP-CLIENT-READ] msgType=0x%02x payloadLen=%d responseFlag=%v\n",
			byte(msgType), len(payload), msgType&MsgResponseFlag != 0)

		// Extract request ID from response
		if len(payload) < 4 {
			continue
		}
		reqID := binary.BigEndian.Uint32(payload[:4])
		respPayload := payload[4:]

		c.mu.Lock()
		ch, ok := c.requests[reqID]
		c.mu.Unlock()

		if ok {
			// Check for error response (MsgResponseFlag set)
			if msgType&MsgResponseFlag != 0 {
				errMsg := string(respPayload)
				ch <- &response{
					err: fmt.Errorf("remote error: %s", errMsg),
				}
			} else {
				ch <- &response{
					msgType: msgType,
					payload: respPayload,
				}
			}
		}
	}
}

// Listener accepts ZAP connections
type Listener struct {
	listener net.Listener
	config   *Config
}

// Listen creates a new ZAP listener
func Listen(addr string, config *Config) (*Listener, error) {
	if config == nil {
		config = DefaultConfig()
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("zap: listen failed: %w", err)
	}

	return &Listener{
		listener: listener,
		config:   config,
	}, nil
}

// Accept accepts a new connection
func (l *Listener) Accept() (*ServerConn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	return newServerConn(conn, l.config), nil
}

// Close closes the listener
func (l *Listener) Close() error {
	return l.listener.Close()
}

// Addr returns the listener's address
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

// ServerConn is a server-side ZAP connection
type ServerConn struct {
	conn   net.Conn
	config *Config

	reader *bufio.Reader
	writer *bufio.Writer

	writeMu sync.Mutex
	closed  atomic.Bool
}

func newServerConn(conn net.Conn, config *Config) *ServerConn {
	return &ServerConn{
		conn:   conn,
		config: config,
		reader: bufio.NewReaderSize(conn, config.BufferSize),
		writer: bufio.NewWriterSize(conn, config.BufferSize),
	}
}

// Read reads the next request from the connection
func (c *ServerConn) Read() (uint32, MessageType, []byte, error) {
	if c.closed.Load() {
		return 0, 0, nil, ErrClosed
	}

	if c.config.ReadTimeout > 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
	}

	msgType, payload, err := ReadMessage(c.reader)
	if err != nil {
		return 0, 0, nil, err
	}

	// Extract request ID
	if len(payload) < 4 {
		return 0, 0, nil, ErrInvalidMessage
	}
	reqID := binary.BigEndian.Uint32(payload[:4])

	return reqID, msgType, payload[4:], nil
}

// Write writes a response
func (c *ServerConn) Write(reqID uint32, msgType MessageType, payload []byte) error {
	if c.closed.Load() {
		return ErrClosed
	}

	// Prepend request ID to response
	respBuf := GetBuffer()
	defer PutBuffer(respBuf)

	respBuf.WriteUint32(reqID)
	respBuf.Grow(respBuf.Len() + len(payload))
	copy(respBuf.Data[respBuf.offset:], payload)
	respBuf.offset += len(payload)

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	fmt.Fprintf(os.Stderr, "[ZAP-SERVER-WRITE] msgType=0x%02x reqID=%d payloadLen=%d responseFlag=%v\n",
		byte(msgType), reqID, len(payload), msgType&MsgResponseFlag != 0)

	if c.config.WriteTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
	}

	if err := WriteMessage(c.writer, msgType, respBuf.Bytes()); err != nil {
		return err
	}
	return c.writer.Flush()
}

// Close closes the server connection
func (c *ServerConn) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	return c.conn.Close()
}

// RemoteAddr returns the remote address
func (c *ServerConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Handler processes ZAP requests
type Handler interface {
	// Handle processes a request and returns a response
	Handle(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error)
}

// HandlerFunc is a function that implements Handler
type HandlerFunc func(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error)

// Handle implements Handler
func (f HandlerFunc) Handle(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error) {
	return f(ctx, msgType, payload)
}

// Server serves ZAP requests
type Server struct {
	listener *Listener
	handler  Handler

	mu     sync.Mutex
	conns  map[*ServerConn]struct{}
	closed atomic.Bool
	done   chan struct{}
}

// NewServer creates a new ZAP server
func NewServer(listener *Listener, handler Handler) *Server {
	return &Server{
		listener: listener,
		handler:  handler,
		conns:    make(map[*ServerConn]struct{}),
		done:     make(chan struct{}),
	}
}

// Serve accepts and processes connections
func (s *Server) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.done:
			return nil
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return nil
			}
			continue
		}

		s.mu.Lock()
		s.conns[conn] = struct{}{}
		s.mu.Unlock()

		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn *ServerConn) {
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.conns, conn)
		s.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		default:
		}

		reqID, msgType, payload, err := conn.Read()
		if err != nil {
			if err == io.EOF || errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}

		// Handle request concurrently to avoid blocking on long-running
		// operations like WaitForEvent. Write is protected by writeMu.
		go func(reqID uint32, msgType MessageType, payload []byte) {
			defer func() {
				if r := recover(); r != nil {
					// Send panic as error response
					errMsg := fmt.Sprintf("handler panic: %v", r)
					conn.Write(reqID, msgType|MsgResponseFlag, []byte(errMsg))
				}
			}()

			respType, respPayload, err := s.handler.Handle(ctx, msgType, payload)
			fmt.Fprintf(os.Stderr, "[ZAP-HANDLER] reqMsgType=0x%02x respType=0x%02x err=%v respLen=%d\n",
				byte(msgType), byte(respType), err, len(respPayload))
			if err != nil {
				conn.Write(reqID, msgType|MsgResponseFlag, []byte(err.Error()))
				return
			}

			conn.Write(reqID, respType, respPayload)
		}(reqID, msgType, payload)
	}
}

// Close closes the server
func (s *Server) Close() error {
	if s.closed.Swap(true) {
		return nil
	}
	close(s.done)

	s.mu.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.mu.Unlock()

	return s.listener.Close()
}
