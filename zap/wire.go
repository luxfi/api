// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
)

const (
	// MaxMessageSize is the maximum allowed message size (16MB)
	MaxMessageSize = 16 * 1024 * 1024

	// HeaderSize is the size of the message header (4 bytes length + 1 byte type)
	HeaderSize = 5

	// DefaultBufferSize for pooled buffers
	DefaultBufferSize = 64 * 1024
)

var (
	ErrMessageTooLarge = errors.New("zap: message exceeds maximum size")
	ErrInvalidMessage  = errors.New("zap: invalid message format")
	ErrUnknownType     = errors.New("zap: unknown message type")
)

// MessageType identifies the RPC method being called
type MessageType uint8

const (
	// VM service methods (1-31)
	MsgInitialize MessageType = iota + 1
	MsgSetState
	MsgShutdown
	MsgCreateHandlers
	MsgNewHTTPHandler
	MsgWaitForEvent
	MsgConnected
	MsgDisconnected
	MsgBuildBlock
	MsgParseBlock
	MsgGetBlock
	MsgSetPreference
	MsgHealth
	MsgVersion
	MsgRequest
	MsgRequestFailed
	MsgResponse
	MsgGossip
	MsgGather
	MsgGetAncestors
	MsgBatchedParseBlock
	MsgGetBlockIDAtHeight
	MsgStateSyncEnabled
	MsgGetOngoingSyncStateSummary
	MsgGetLastStateSummary
	MsgParseStateSummary
	MsgGetStateSummary
	MsgBlockVerify
	MsgBlockAccept
	MsgBlockReject
	MsgStateSummaryAccept // = 31

	// p2p.Sender methods (40-49)
	MsgSendRequest  MessageType = 40
	MsgSendResponse MessageType = 41
	MsgSendError    MessageType = 42
	MsgSendGossip   MessageType = 43

	// Warp signing methods (50-59)
	MsgWarpSign         MessageType = 50
	MsgWarpGetPublicKey MessageType = 51
	MsgWarpBatchSign    MessageType = 52

	// Response flag - set on response messages (high bit)
	// All message types must be < 128 to allow OR with this flag
	MsgResponseFlag MessageType = 128
)

// Buffer is a reusable byte buffer for zero-copy operations
type Buffer struct {
	Data   []byte
	offset int
}

// BufferPool manages reusable buffers to minimize allocations
var BufferPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{
			Data: make([]byte, DefaultBufferSize),
		}
	},
}

// GetBuffer retrieves a buffer from the pool
func GetBuffer() *Buffer {
	buf := BufferPool.Get().(*Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *Buffer) {
	if buf != nil && cap(buf.Data) <= MaxMessageSize {
		BufferPool.Put(buf)
	}
}

// Reset prepares the buffer for reuse
func (b *Buffer) Reset() {
	b.offset = 0
	b.Data = b.Data[:cap(b.Data)]
}

// Grow ensures the buffer has at least n bytes available
func (b *Buffer) Grow(n int) {
	if cap(b.Data) < n {
		newData := make([]byte, n*2)
		copy(newData, b.Data[:b.offset])
		b.Data = newData
	}
	if len(b.Data) < n {
		b.Data = b.Data[:cap(b.Data)]
	}
}

// WriteUint8 writes a uint8 to the buffer
func (b *Buffer) WriteUint8(v uint8) {
	b.Grow(b.offset + 1)
	b.Data[b.offset] = v
	b.offset++
}

// WriteUint16 writes a uint16 to the buffer (big-endian)
func (b *Buffer) WriteUint16(v uint16) {
	b.Grow(b.offset + 2)
	binary.BigEndian.PutUint16(b.Data[b.offset:], v)
	b.offset += 2
}

// WriteUint32 writes a uint32 to the buffer (big-endian)
func (b *Buffer) WriteUint32(v uint32) {
	b.Grow(b.offset + 4)
	binary.BigEndian.PutUint32(b.Data[b.offset:], v)
	b.offset += 4
}

// WriteUint64 writes a uint64 to the buffer (big-endian)
func (b *Buffer) WriteUint64(v uint64) {
	b.Grow(b.offset + 8)
	binary.BigEndian.PutUint64(b.Data[b.offset:], v)
	b.offset += 8
}

// WriteInt32 writes an int32 to the buffer (big-endian)
func (b *Buffer) WriteInt32(v int32) {
	b.WriteUint32(uint32(v))
}

// WriteInt64 writes an int64 to the buffer (big-endian)
func (b *Buffer) WriteInt64(v int64) {
	b.WriteUint64(uint64(v))
}

// WriteBool writes a boolean to the buffer
func (b *Buffer) WriteBool(v bool) {
	if v {
		b.WriteUint8(1)
	} else {
		b.WriteUint8(0)
	}
}

// WriteBytes writes a length-prefixed byte slice to the buffer
func (b *Buffer) WriteBytes(data []byte) {
	b.WriteUint32(uint32(len(data)))
	b.Grow(b.offset + len(data))
	copy(b.Data[b.offset:], data)
	b.offset += len(data)
}

// WriteString writes a length-prefixed string to the buffer
func (b *Buffer) WriteString(s string) {
	b.WriteBytes([]byte(s))
}

// Bytes returns the written portion of the buffer
func (b *Buffer) Bytes() []byte {
	return b.Data[:b.offset]
}

// Len returns the number of bytes written
func (b *Buffer) Len() int {
	return b.offset
}

// Reader provides zero-copy reading from a byte slice
type Reader struct {
	data   []byte
	offset int
}

// NewReader creates a new reader from a byte slice
func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

// ReadUint8 reads a uint8 from the buffer
func (r *Reader) ReadUint8() (uint8, error) {
	if r.offset+1 > len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := r.data[r.offset]
	r.offset++
	return v, nil
}

// ReadUint16 reads a uint16 from the buffer (big-endian)
func (r *Reader) ReadUint16() (uint16, error) {
	if r.offset+2 > len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := binary.BigEndian.Uint16(r.data[r.offset:])
	r.offset += 2
	return v, nil
}

// ReadUint32 reads a uint32 from the buffer (big-endian)
func (r *Reader) ReadUint32() (uint32, error) {
	if r.offset+4 > len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := binary.BigEndian.Uint32(r.data[r.offset:])
	r.offset += 4
	return v, nil
}

// ReadUint64 reads a uint64 from the buffer (big-endian)
func (r *Reader) ReadUint64() (uint64, error) {
	if r.offset+8 > len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := binary.BigEndian.Uint64(r.data[r.offset:])
	r.offset += 8
	return v, nil
}

// ReadInt32 reads an int32 from the buffer (big-endian)
func (r *Reader) ReadInt32() (int32, error) {
	v, err := r.ReadUint32()
	return int32(v), err
}

// ReadInt64 reads an int64 from the buffer (big-endian)
func (r *Reader) ReadInt64() (int64, error) {
	v, err := r.ReadUint64()
	return int64(v), err
}

// ReadBool reads a boolean from the buffer
func (r *Reader) ReadBool() (bool, error) {
	v, err := r.ReadUint8()
	return v != 0, err
}

// ReadBytes reads a length-prefixed byte slice (zero-copy - returns slice into original buffer)
func (r *Reader) ReadBytes() ([]byte, error) {
	length, err := r.ReadUint32()
	if err != nil {
		return nil, err
	}
	if r.offset+int(length) > len(r.data) {
		return nil, io.ErrUnexpectedEOF
	}
	// Zero-copy: return slice into original buffer
	data := r.data[r.offset : r.offset+int(length)]
	r.offset += int(length)
	return data, nil
}

// ReadString reads a length-prefixed string
func (r *Reader) ReadString() (string, error) {
	data, err := r.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Remaining returns the number of unread bytes
func (r *Reader) Remaining() int {
	return len(r.data) - r.offset
}

// WriteMessage writes a complete ZAP message with header
func WriteMessage(w io.Writer, msgType MessageType, payload []byte) error {
	if len(payload) > MaxMessageSize {
		return ErrMessageTooLarge
	}

	header := make([]byte, HeaderSize)
	binary.BigEndian.PutUint32(header[0:4], uint32(len(payload)))
	header[4] = byte(msgType)

	if _, err := w.Write(header); err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}
	return nil
}

// ReadMessage reads a complete ZAP message with header
func ReadMessage(r io.Reader) (MessageType, []byte, error) {
	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, err
	}

	length := binary.BigEndian.Uint32(header[0:4])
	msgType := MessageType(header[4])

	if length > MaxMessageSize {
		return 0, nil, ErrMessageTooLarge
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, err
	}

	return msgType, payload, nil
}
