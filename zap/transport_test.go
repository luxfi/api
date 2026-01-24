// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import (
	"bytes"
	"context"
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
