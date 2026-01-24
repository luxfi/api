// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	// Write various types
	buf.WriteUint8(0x42)
	buf.WriteUint16(0x1234)
	buf.WriteUint32(0xDEADBEEF)
	buf.WriteUint64(0xCAFEBABEDEADBEEF)
	buf.WriteInt32(-12345)
	buf.WriteInt64(-9876543210)
	buf.WriteBool(true)
	buf.WriteBool(false)
	buf.WriteBytes([]byte("hello"))
	buf.WriteString("world")

	// Read back
	r := NewReader(buf.Bytes())

	v8, err := r.ReadUint8()
	if err != nil || v8 != 0x42 {
		t.Errorf("ReadUint8: got %v, %v", v8, err)
	}

	v16, err := r.ReadUint16()
	if err != nil || v16 != 0x1234 {
		t.Errorf("ReadUint16: got %v, %v", v16, err)
	}

	v32, err := r.ReadUint32()
	if err != nil || v32 != 0xDEADBEEF {
		t.Errorf("ReadUint32: got %v, %v", v32, err)
	}

	v64, err := r.ReadUint64()
	if err != nil || v64 != 0xCAFEBABEDEADBEEF {
		t.Errorf("ReadUint64: got %v, %v", v64, err)
	}

	i32, err := r.ReadInt32()
	if err != nil || i32 != -12345 {
		t.Errorf("ReadInt32: got %v, %v", i32, err)
	}

	i64, err := r.ReadInt64()
	if err != nil || i64 != -9876543210 {
		t.Errorf("ReadInt64: got %v, %v", i64, err)
	}

	b1, err := r.ReadBool()
	if err != nil || !b1 {
		t.Errorf("ReadBool: got %v, %v", b1, err)
	}

	b2, err := r.ReadBool()
	if err != nil || b2 {
		t.Errorf("ReadBool: got %v, %v", b2, err)
	}

	data, err := r.ReadBytes()
	if err != nil || string(data) != "hello" {
		t.Errorf("ReadBytes: got %v, %v", string(data), err)
	}

	str, err := r.ReadString()
	if err != nil || str != "world" {
		t.Errorf("ReadString: got %v, %v", str, err)
	}

	if r.Remaining() != 0 {
		t.Errorf("Expected 0 remaining, got %d", r.Remaining())
	}
}

func TestMessage(t *testing.T) {
	var buf bytes.Buffer

	// Write message
	payload := []byte("test payload")
	if err := WriteMessage(&buf, MsgBuildBlock, payload); err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}

	// Read message
	msgType, data, err := ReadMessage(&buf)
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}

	if msgType != MsgBuildBlock {
		t.Errorf("Expected MsgBuildBlock, got %d", msgType)
	}

	if !bytes.Equal(data, payload) {
		t.Errorf("Payload mismatch: got %v, want %v", data, payload)
	}
}

func TestBlockResponse(t *testing.T) {
	orig := &BlockResponse{
		ID:                []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		ParentID:          []byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		Bytes:             make([]byte, 1024),
		Height:            12345,
		Timestamp:         1704067200,
		VerifyWithContext: true,
		Err:               ErrorUnspecified,
	}

	// Fill block bytes with test data
	for i := range orig.Bytes {
		orig.Bytes[i] = byte(i % 256)
	}

	buf := GetBuffer()
	defer PutBuffer(buf)

	orig.Encode(buf)

	decoded := &BlockResponse{}
	if err := decoded.Decode(NewReader(buf.Bytes())); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if !bytes.Equal(decoded.ID, orig.ID) {
		t.Error("ID mismatch")
	}
	if !bytes.Equal(decoded.ParentID, orig.ParentID) {
		t.Error("ParentID mismatch")
	}
	if !bytes.Equal(decoded.Bytes, orig.Bytes) {
		t.Error("Bytes mismatch")
	}
	if decoded.Height != orig.Height {
		t.Errorf("Height: got %d, want %d", decoded.Height, orig.Height)
	}
	if decoded.Timestamp != orig.Timestamp {
		t.Errorf("Timestamp: got %d, want %d", decoded.Timestamp, orig.Timestamp)
	}
	if decoded.VerifyWithContext != orig.VerifyWithContext {
		t.Error("VerifyWithContext mismatch")
	}
	if decoded.Err != orig.Err {
		t.Error("Err mismatch")
	}
}

func TestBatchedParseBlockRequest(t *testing.T) {
	orig := &BatchedParseBlockRequest{
		Requests: [][]byte{
			make([]byte, 100),
			make([]byte, 200),
			make([]byte, 300),
		},
	}

	for i, req := range orig.Requests {
		for j := range req {
			req[j] = byte((i + j) % 256)
		}
	}

	buf := GetBuffer()
	defer PutBuffer(buf)

	orig.Encode(buf)

	decoded := &BatchedParseBlockRequest{}
	if err := decoded.Decode(NewReader(buf.Bytes())); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if len(decoded.Requests) != len(orig.Requests) {
		t.Fatalf("Request count: got %d, want %d", len(decoded.Requests), len(orig.Requests))
	}

	for i := range orig.Requests {
		if !bytes.Equal(decoded.Requests[i], orig.Requests[i]) {
			t.Errorf("Request %d mismatch", i)
		}
	}
}

func TestSendGossipMsg(t *testing.T) {
	orig := &SendGossipMsg{
		NodeIDs: [][]byte{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		},
		Validators:    10,
		NonValidators: 5,
		Peers:         100,
		Msg:           []byte("gossip message payload"),
	}

	buf := GetBuffer()
	defer PutBuffer(buf)

	orig.Encode(buf)

	decoded := &SendGossipMsg{}
	if err := decoded.Decode(NewReader(buf.Bytes())); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if len(decoded.NodeIDs) != len(orig.NodeIDs) {
		t.Fatalf("NodeID count: got %d, want %d", len(decoded.NodeIDs), len(orig.NodeIDs))
	}

	for i := range orig.NodeIDs {
		if !bytes.Equal(decoded.NodeIDs[i], orig.NodeIDs[i]) {
			t.Errorf("NodeID %d mismatch", i)
		}
	}

	if decoded.Validators != orig.Validators {
		t.Errorf("Validators: got %d, want %d", decoded.Validators, orig.Validators)
	}
	if decoded.NonValidators != orig.NonValidators {
		t.Errorf("NonValidators: got %d, want %d", decoded.NonValidators, orig.NonValidators)
	}
	if decoded.Peers != orig.Peers {
		t.Errorf("Peers: got %d, want %d", decoded.Peers, orig.Peers)
	}
	if !bytes.Equal(decoded.Msg, orig.Msg) {
		t.Error("Msg mismatch")
	}
}

func BenchmarkBufferWrite(b *testing.B) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	data := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteUint64(12345678901234567890)
		buf.WriteBytes(data)
		buf.WriteString("test string")
	}
}

func BenchmarkBufferRead(b *testing.B) {
	buf := GetBuffer()
	buf.WriteUint64(12345678901234567890)
	buf.WriteBytes(make([]byte, 1024))
	buf.WriteString("test string")
	data := buf.Bytes()
	PutBuffer(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(data)
		_, _ = r.ReadUint64()
		_, _ = r.ReadBytes()
		_, _ = r.ReadString()
	}
}

func BenchmarkBlockResponse(b *testing.B) {
	resp := &BlockResponse{
		ID:                make([]byte, 32),
		ParentID:          make([]byte, 32),
		Bytes:             make([]byte, 100*1024), // 100KB block
		Height:            12345,
		Timestamp:         1704067200,
		VerifyWithContext: true,
		Err:               ErrorUnspecified,
	}

	buf := GetBuffer()
	defer PutBuffer(buf)

	b.Run("Encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf.Reset()
			resp.Encode(buf)
		}
		b.SetBytes(int64(buf.Len()))
	})

	resp.Encode(buf)
	data := make([]byte, len(buf.Bytes()))
	copy(data, buf.Bytes())

	b.Run("Decode", func(b *testing.B) {
		decoded := &BlockResponse{}
		for i := 0; i < b.N; i++ {
			r := NewReader(data)
			_ = decoded.Decode(r)
		}
		b.SetBytes(int64(len(data)))
	})
}
