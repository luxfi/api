// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import "testing"

// TestInitializeResponseDecodeLengthTolerant proves the OPTIONAL trailing
// Capabilities field can be absent on the wire without breaking Decode — the
// forward/backward-compatible property that fixes the v1.36.11 "unexpected EOF"
// skew, where a node built with the Capabilities field decoded an InitializeResponse
// from a stale plugin that never encoded it.
func TestInitializeResponseDecodeLengthTolerant(t *testing.T) {
	lastID := []byte("last-accepted-id-32-bytes-000000")
	parentID := []byte("parent-accepted-id-32-bytes-0000")
	const height = uint64(12345)
	blockBytes := []byte("block-bytes")
	const ts = int64(1_700_000_000_000_000_000)

	// PRE-v1.0.16 payload: the five original fields only, exactly as a stale
	// plugin (built before Capabilities existed) encodes it.
	old := GetBuffer()
	old.WriteBytes(lastID)
	old.WriteBytes(parentID)
	old.WriteUint64(height)
	old.WriteBytes(blockBytes)
	old.WriteInt64(ts)
	oldPayload := append([]byte(nil), old.Bytes()...)
	PutBuffer(old)

	var got InitializeResponse
	if err := got.Decode(NewReader(oldPayload)); err != nil {
		t.Fatalf("decode of pre-Capabilities payload must not error, got: %v", err)
	}
	if got.Capabilities != 0 {
		t.Fatalf("absent Capabilities must default to 0, got %d", got.Capabilities)
	}
	if string(got.LastAcceptedID) != string(lastID) ||
		string(got.LastAcceptedParentID) != string(parentID) ||
		got.Height != height ||
		string(got.Bytes) != string(blockBytes) ||
		got.Timestamp != ts {
		t.Fatalf("non-optional fields corrupted by tolerant decode: %+v", got)
	}

	// FULL round-trip: Capabilities present survives Encode -> Decode.
	want := InitializeResponse{
		LastAcceptedID:       lastID,
		LastAcceptedParentID: parentID,
		Height:               height,
		Bytes:                blockBytes,
		Timestamp:            ts,
		Capabilities:         CapQuasarExport,
	}
	buf := GetBuffer()
	want.Encode(buf)
	full := append([]byte(nil), buf.Bytes()...)
	PutBuffer(buf)

	var rt InitializeResponse
	if err := rt.Decode(NewReader(full)); err != nil {
		t.Fatalf("round-trip decode: %v", err)
	}
	if rt.Capabilities != CapQuasarExport {
		t.Fatalf("Capabilities round-trip: got %d want %d", rt.Capabilities, CapQuasarExport)
	}

	// FORWARD-compat: an even-newer peer appends further trailing bytes; Decode
	// reads Capabilities and ignores the rest without erroring.
	fwd := append(append([]byte(nil), full...), 0xDE, 0xAD, 0xBE, 0xEF)
	var fr InitializeResponse
	if err := fr.Decode(NewReader(fwd)); err != nil {
		t.Fatalf("forward-compat decode must ignore trailing bytes, got: %v", err)
	}
	if fr.Capabilities != CapQuasarExport {
		t.Fatalf("forward-compat Capabilities: got %d want %d", fr.Capabilities, CapQuasarExport)
	}
}
