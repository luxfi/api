// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import (
	"bytes"
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
	// Start server
	listener, err := Listen("127.0.0.1:0", nil)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer listener.Close()

	// Echo handler - returns the same message type and payload
	handler := HandlerFunc(func(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error) {
		return msgType, payload, nil
	})

	server := NewServer(listener, handler)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = server.Serve(ctx)
	}()

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Connect client
	conn, err := Dial(ctx, listener.Addr().String(), nil)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	// Test request/response
	payload := []byte("test payload")
	respType, respPayload, err := conn.Call(ctx, MsgBuildBlock, payload)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}

	if respType != MsgBuildBlock|MsgResponseFlag {
		t.Errorf("Expected response type %d, got %d", MsgBuildBlock|MsgResponseFlag, respType)
	}

	if !bytes.Equal(respPayload, payload) {
		t.Errorf("Payload mismatch: got %v, want %v", respPayload, payload)
	}

	// Clean up
	cancel()
	server.Close()
	wg.Wait()
}

func TestConcurrentCalls(t *testing.T) {
	// Start server
	listener, err := Listen("127.0.0.1:0", nil)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer listener.Close()

	// Echo handler with small delay
	handler := HandlerFunc(func(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error) {
		return msgType, payload, nil
	})

	server := NewServer(listener, handler)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = server.Serve(ctx)
	}()

	time.Sleep(10 * time.Millisecond)

	// Connect client
	conn, err := Dial(ctx, listener.Addr().String(), nil)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	// Launch concurrent calls
	const numCalls = 100
	var wg sync.WaitGroup
	errors := make(chan error, numCalls)

	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := []byte{byte(i)}
			_, respPayload, err := conn.Call(ctx, MsgParseBlock, payload)
			if err != nil {
				errors <- err
				return
			}
			if !bytes.Equal(respPayload, payload) {
				errors <- ErrResponseFailed
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent call error: %v", err)
	}

	cancel()
	server.Close()
}

func BenchmarkRoundTrip(b *testing.B) {
	listener, err := Listen("127.0.0.1:0", nil)
	if err != nil {
		b.Fatalf("Listen: %v", err)
	}
	defer listener.Close()

	handler := HandlerFunc(func(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error) {
		return msgType, payload, nil
	})

	server := NewServer(listener, handler)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = server.Serve(ctx)
	}()

	time.Sleep(10 * time.Millisecond)

	conn, err := Dial(ctx, listener.Addr().String(), nil)
	if err != nil {
		b.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	payload := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := conn.Call(ctx, MsgBuildBlock, payload)
		if err != nil {
			b.Fatalf("Call: %v", err)
		}
	}
	b.SetBytes(int64(len(payload)))

	cancel()
	server.Close()
}

func BenchmarkRoundTripLargeBlock(b *testing.B) {
	listener, err := Listen("127.0.0.1:0", nil)
	if err != nil {
		b.Fatalf("Listen: %v", err)
	}
	defer listener.Close()

	handler := HandlerFunc(func(ctx context.Context, msgType MessageType, payload []byte) (MessageType, []byte, error) {
		return msgType, payload, nil
	})

	server := NewServer(listener, handler)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = server.Serve(ctx)
	}()

	time.Sleep(10 * time.Millisecond)

	conn, err := Dial(ctx, listener.Addr().String(), nil)
	if err != nil {
		b.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	// Simulate 100KB block
	payload := make([]byte, 100*1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := conn.Call(ctx, MsgBuildBlock, payload)
		if err != nil {
			b.Fatalf("Call: %v", err)
		}
	}
	b.SetBytes(int64(len(payload)))

	cancel()
	server.Close()
}

// gatedListener hands out exactly one connection, and only after release is
// closed. It makes the Close/Accept interleaving deterministic instead of
// hoping a real socket race reproduces.
type gatedListener struct {
	entered  chan struct{} // closed once Serve is parked inside Accept
	release  chan struct{}
	accepted bool
	conn     net.Conn
}

func (l *gatedListener) Accept() (net.Conn, error) {
	if l.accepted {
		return nil, net.ErrClosed
	}
	close(l.entered)
	<-l.release
	l.accepted = true
	return l.conn, nil
}
func (l *gatedListener) Close() error   { return nil }
func (l *gatedListener) Addr() net.Addr { return l.conn.LocalAddr() }

// TestServer_AcceptAfterCloseDoesNotPanic pins a shutdown race that killed the
// whole process. Close() nils the conns map under s.mu; Serve's accept loop then
// assigned into that nil map and panicked. Serve runs on a detached goroutine,
// so the panic was fatal rather than returned to anyone. Every server whose stop
// path is "cancel then Close" while a connection is in flight could hit it —
// which includes the node's per-chain rpcdb and atomic shared-memory servers.
//
// The interleaving is forced, not raced: Serve is parked inside Accept, Close
// runs to completion, and only then does Accept yield a connection.
func TestServer_AcceptAfterCloseDoesNotPanic(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	l := &gatedListener{entered: make(chan struct{}), release: make(chan struct{}), conn: server}
	srv := NewServer(NewListener(l, DefaultConfig()), HandlerFunc(
		func(_ context.Context, msgType MessageType, _ []byte) (MessageType, []byte, error) {
			return msgType, nil, nil
		}))

	served := make(chan error, 1)
	go func() { served <- srv.Serve(context.Background()) }()

	// Wait until Serve is genuinely parked inside Accept — otherwise it may see
	// the closed done channel at the top of the loop and return before ever
	// reaching the assignment this test is about.
	<-l.entered
	// Close first: conns becomes nil.
	srv.Close()
	// Now let Accept deliver a connection into the closed server.
	close(l.release)

	select {
	case err := <-served:
		if err != nil {
			t.Fatalf("Serve returned %v, want nil after Close", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return after Close")
	}
}
