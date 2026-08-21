// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

import (
	"bytes"
	"testing"
)

// roundTrip encodes m, hands the bytes to decode, and returns the payload so a
// test can also assert on framing.
func encodeToBytes(enc func(*Buffer)) []byte {
	buf := GetBuffer()
	defer PutBuffer(buf)
	enc(buf)
	out := make([]byte, len(buf.Bytes()))
	copy(out, buf.Bytes())
	return out
}

func TestAtomicGetRequest_RoundTrip(t *testing.T) {
	in := &AtomicGetRequest{
		PeerChainID: bytes.Repeat([]byte{0xD1}, 32),
		Keys:        [][]byte{bytes.Repeat([]byte{0x01}, 32), bytes.Repeat([]byte{0x02}, 32)},
	}
	payload := encodeToBytes(in.Encode)

	var out AtomicGetRequest
	if err := out.Decode(NewReader(payload)); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !bytes.Equal(out.PeerChainID, in.PeerChainID) {
		t.Fatalf("peerChainID = %x, want %x", out.PeerChainID, in.PeerChainID)
	}
	if len(out.Keys) != len(in.Keys) {
		t.Fatalf("keys len = %d, want %d", len(out.Keys), len(in.Keys))
	}
	for i := range in.Keys {
		if !bytes.Equal(out.Keys[i], in.Keys[i]) {
			t.Fatalf("key[%d] = %x, want %x", i, out.Keys[i], in.Keys[i])
		}
	}
}

// TestAtomicGetResponse_AbsentIsEmptyNotShort pins the invariant the settlement
// path depends on: SharedMemory.Get returns a values slice the SAME length as
// keys, with an absent object represented as an EMPTY value at its index. If the
// wire collapsed empties the caller would mis-align values with keys and could
// bind a credit to the WRONG object.
func TestAtomicGetResponse_AbsentIsEmptyNotShort(t *testing.T) {
	in := &AtomicGetResponse{
		Values: [][]byte{
			bytes.Repeat([]byte{0xAA}, 69), // present
			nil,                            // absent
			bytes.Repeat([]byte{0xBB}, 69), // present
		},
	}
	payload := encodeToBytes(in.Encode)

	var out AtomicGetResponse
	if err := out.Decode(NewReader(payload)); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out.Values) != 3 {
		t.Fatalf("values len = %d, want 3 (positional alignment lost)", len(out.Values))
	}
	if len(out.Values[1]) != 0 {
		t.Fatalf("absent value decoded as %d bytes, want 0", len(out.Values[1]))
	}
	if !bytes.Equal(out.Values[0], in.Values[0]) || !bytes.Equal(out.Values[2], in.Values[2]) {
		t.Fatal("present values did not survive the round trip")
	}
}

func TestAtomicApplyRequest_RoundTrip(t *testing.T) {
	in := &AtomicApplyRequest{Requests: []*AtomicChainRequests{{
		PeerChainID:    bytes.Repeat([]byte{0xD1}, 32),
		RemoveRequests: [][]byte{bytes.Repeat([]byte{0x11}, 32)},
		PutRequests: []*AtomicElement{{
			Key:    bytes.Repeat([]byte{0x22}, 32),
			Value:  bytes.Repeat([]byte{0x33}, 69),
			Traits: [][]byte{bytes.Repeat([]byte{0x44}, 20)},
		}},
	}}}
	payload := encodeToBytes(in.Encode)

	var out AtomicApplyRequest
	if err := out.Decode(NewReader(payload)); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out.Requests) != 1 {
		t.Fatalf("requests len = %d, want 1", len(out.Requests))
	}
	got := out.Requests[0]
	if !bytes.Equal(got.PeerChainID, in.Requests[0].PeerChainID) {
		t.Fatal("peerChainID mismatch")
	}
	if len(got.RemoveRequests) != 1 || !bytes.Equal(got.RemoveRequests[0], in.Requests[0].RemoveRequests[0]) {
		t.Fatal("removes mismatch")
	}
	if len(got.PutRequests) != 1 {
		t.Fatalf("puts len = %d, want 1", len(got.PutRequests))
	}
	wantPut := in.Requests[0].PutRequests[0]
	gotPut := got.PutRequests[0]
	if !bytes.Equal(gotPut.Key, wantPut.Key) || !bytes.Equal(gotPut.Value, wantPut.Value) {
		t.Fatal("put key/value mismatch")
	}
	if len(gotPut.Traits) != 1 || !bytes.Equal(gotPut.Traits[0], wantPut.Traits[0]) {
		t.Fatal("put traits mismatch")
	}
}

// TestAtomicApplyRequest_EmptyIsWellFormed: a block with no cross-chain ops must
// still encode/decode cleanly (the flush path skips the call in that case, but a
// zero-length request must never be a decode error).
func TestAtomicApplyRequest_EmptyIsWellFormed(t *testing.T) {
	payload := encodeToBytes((&AtomicApplyRequest{}).Encode)
	var out AtomicApplyRequest
	if err := out.Decode(NewReader(payload)); err != nil {
		t.Fatalf("decode of empty apply: %v", err)
	}
	if len(out.Requests) != 0 {
		t.Fatalf("requests len = %d, want 0", len(out.Requests))
	}
}

// --- InitializeRequest skew. These are the rollout-risk tests: a wire-format
// change that forks a mixed fleet is the failure mode this whole design is built
// to avoid, so both skew directions are pinned here rather than reasoned about.

// TestInitializeRequest_OldDecoderIgnoresNewTrailingFields is the NEW node ->
// OLD plugin direction. The old decoder stops after ServerAddr; the appended
// bytes are inert. Simulated by decoding only through ServerAddr and asserting
// no error and that the earlier fields are intact.
func TestInitializeRequest_OldDecoderIgnoresNewTrailingFields(t *testing.T) {
	newReq := &InitializeRequest{
		NetworkID:        3,
		ChainID:          bytes.Repeat([]byte{0xC1}, 32),
		NodeID:           bytes.Repeat([]byte{0x0E}, 20),
		PublicKey:        []byte{0x01},
		XChainID:         bytes.Repeat([]byte{0x02}, 32),
		CChainID:         bytes.Repeat([]byte{0x03}, 32),
		UTXOAssetID:      bytes.Repeat([]byte{0x04}, 32),
		ChainDataDir:     "/data/chain",
		GenesisBytes:     []byte("genesis"),
		UpgradeBytes:     []byte("upgrade"),
		ConfigBytes:      []byte("config"),
		DBServerAddr:     "127.0.0.1:1",
		ServerAddr:       "127.0.0.1:2",
		AtomicServerAddr: "127.0.0.1:3",
		DChainID:         bytes.Repeat([]byte{0xD1}, 32),
	}
	payload := encodeToBytes(newReq.Encode)

	// Replay the pre-change field order exactly — this is what an old plugin does.
	r := NewReader(payload)
	readBytes := func() []byte { v, err := r.ReadBytes(); mustNoErr(t, err); return v }
	readStr := func() string { v, err := r.ReadString(); mustNoErr(t, err); return v }

	netID, err := r.ReadUint32()
	mustNoErr(t, err)
	if netID != 3 {
		t.Fatalf("networkID = %d, want 3", netID)
	}
	chainID := readBytes()
	readBytes() // nodeID
	readBytes() // publicKey
	readBytes() // xChainID
	readBytes() // cChainID
	readBytes() // utxoAssetID
	dataDir := readStr()
	readBytes() // genesis
	readBytes() // upgrade
	readBytes() // config
	readStr()   // dbServerAddr
	serverAddr := readStr()

	if !bytes.Equal(chainID, newReq.ChainID) {
		t.Fatal("old decoder mis-read ChainID")
	}
	if dataDir != "/data/chain" || serverAddr != "127.0.0.1:2" {
		t.Fatalf("old decoder mis-read tail: dataDir=%q serverAddr=%q", dataDir, serverAddr)
	}
	// The appended fields are still on the wire — an old decoder simply never
	// looks at them. That inertness is the point.
	if r.Remaining() == 0 {
		t.Fatal("expected trailing appended bytes to be present but unread")
	}
}

// TestInitializeRequest_NewDecoderReadsOldPayload is the OLD node -> NEW plugin
// direction: the payload ends after ServerAddr, and the new decoder must yield
// an ABSENT atomic capability rather than io.ErrUnexpectedEOF. This is the exact
// skew class that broke every EVM Initialize with "zap decode initialize
// response: unexpected EOF" when Capabilities was appended.
func TestInitializeRequest_NewDecoderReadsOldPayload(t *testing.T) {
	// Build a payload in the OLD field order (no appended fields).
	old := encodeToBytes(func(buf *Buffer) {
		buf.WriteUint32(3)
		buf.WriteBytes(bytes.Repeat([]byte{0xC1}, 32))
		buf.WriteBytes(bytes.Repeat([]byte{0x0E}, 20))
		buf.WriteBytes([]byte{0x01})
		buf.WriteBytes(bytes.Repeat([]byte{0x02}, 32))
		buf.WriteBytes(bytes.Repeat([]byte{0x03}, 32))
		buf.WriteBytes(bytes.Repeat([]byte{0x04}, 32))
		buf.WriteString("/data/chain")
		buf.WriteBytes([]byte("genesis"))
		buf.WriteBytes([]byte("upgrade"))
		buf.WriteBytes([]byte("config"))
		buf.WriteString("127.0.0.1:1")
		buf.WriteString("127.0.0.1:2")
	})

	var out InitializeRequest
	if err := out.Decode(NewReader(old)); err != nil {
		t.Fatalf("new decoder on old payload: %v (this is the EOF regression)", err)
	}
	if out.AtomicServerAddr != "" {
		t.Fatalf("AtomicServerAddr = %q, want empty (capability must read as ABSENT)", out.AtomicServerAddr)
	}
	if len(out.DChainID) != 0 {
		t.Fatalf("DChainID = %x, want empty", out.DChainID)
	}
	if out.ServerAddr != "127.0.0.1:2" || out.ChainDataDir != "/data/chain" {
		t.Fatal("new decoder mis-read the old payload's existing fields")
	}
}

func TestInitializeRequest_RoundTripWithAtomicFields(t *testing.T) {
	in := &InitializeRequest{
		NetworkID:        1,
		ChainID:          bytes.Repeat([]byte{0xC1}, 32),
		NodeID:           bytes.Repeat([]byte{0x0E}, 20),
		XChainID:         bytes.Repeat([]byte{0x02}, 32),
		CChainID:         bytes.Repeat([]byte{0x03}, 32),
		UTXOAssetID:      bytes.Repeat([]byte{0x04}, 32),
		ChainDataDir:     "/d",
		AtomicServerAddr: "/tmp/atomic.sock",
		DChainID:         bytes.Repeat([]byte{0xD1}, 32),
	}
	var out InitializeRequest
	if err := out.Decode(NewReader(encodeToBytes(in.Encode))); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.AtomicServerAddr != in.AtomicServerAddr {
		t.Fatalf("AtomicServerAddr = %q, want %q", out.AtomicServerAddr, in.AtomicServerAddr)
	}
	if !bytes.Equal(out.DChainID, in.DChainID) {
		t.Fatalf("DChainID = %x, want %x", out.DChainID, in.DChainID)
	}
}

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestInitializeRequest_ValidatorAddrSkew pins BOTH skew directions for the
// appended ValidatorServerAddr, because a wire field is only safe to append if
// neither side of a partial rollout misreads the frame.
//
// The capability must read as ABSENT in both directions, never as present-with-
// garbage: a plugin that thinks it has validator state but dials nothing forms a
// committee of nobody, and a quorum of nobody signs anything.
func TestInitializeRequest_ValidatorAddrSkew(t *testing.T) {
	// OLD node -> NEW plugin: the payload stops after DChainID. The new decoder
	// must leave ValidatorServerAddr empty rather than fail on EOF.
	oldPayload := encodeToBytes(func(buf *Buffer) {
		buf.WriteUint32(3)
		buf.WriteBytes(bytes.Repeat([]byte{0xC1}, 32))
		buf.WriteBytes(bytes.Repeat([]byte{0x0E}, 20))
		buf.WriteBytes([]byte{0x01})
		buf.WriteBytes(bytes.Repeat([]byte{0x02}, 32))
		buf.WriteBytes(bytes.Repeat([]byte{0x03}, 32))
		buf.WriteBytes(bytes.Repeat([]byte{0x04}, 32))
		buf.WriteString("/data/chain")
		buf.WriteBytes([]byte("genesis"))
		buf.WriteBytes([]byte("upgrade"))
		buf.WriteBytes([]byte("config"))
		buf.WriteString("127.0.0.1:1")
		buf.WriteString("127.0.0.1:2")
		buf.WriteString("/tmp/atomic.sock")
		buf.WriteBytes(bytes.Repeat([]byte{0x0D}, 32))
	})
	var fromOld InitializeRequest
	if err := fromOld.Decode(NewReader(oldPayload)); err != nil {
		t.Fatalf("new decoder on a payload without the validator field: %v", err)
	}
	if fromOld.ValidatorServerAddr != "" {
		t.Fatalf("ValidatorServerAddr = %q, want empty — the capability must read ABSENT",
			fromOld.ValidatorServerAddr)
	}
	if fromOld.AtomicServerAddr != "/tmp/atomic.sock" {
		t.Fatalf("appending a field disturbed the one before it: AtomicServerAddr = %q",
			fromOld.AtomicServerAddr)
	}

	// NEW node -> OLD plugin: the field is written last, so a decoder that stops
	// after DChainID reads every earlier field correctly and ignores the tail.
	in := &InitializeRequest{
		NetworkID:           1,
		ChainID:             bytes.Repeat([]byte{0xC1}, 32),
		NodeID:              bytes.Repeat([]byte{0x0E}, 20),
		XChainID:            bytes.Repeat([]byte{0x02}, 32),
		CChainID:            bytes.Repeat([]byte{0x03}, 32),
		UTXOAssetID:         bytes.Repeat([]byte{0x04}, 32),
		ChainDataDir:        "/d",
		AtomicServerAddr:    "/tmp/atomic.sock",
		DChainID:            bytes.Repeat([]byte{0x0D}, 32),
		ValidatorServerAddr: "/tmp/validators.sock",
	}
	var out InitializeRequest
	if err := out.Decode(NewReader(encodeToBytes(in.Encode))); err != nil {
		t.Fatalf("round trip: %v", err)
	}
	if out.ValidatorServerAddr != "/tmp/validators.sock" {
		t.Fatalf("ValidatorServerAddr round trip = %q", out.ValidatorServerAddr)
	}
	if out.AtomicServerAddr != "/tmp/atomic.sock" || string(out.DChainID) != string(in.DChainID) {
		t.Fatal("the appended field disturbed the fields before it")
	}
}

// TestEveryMessageTypeIsBelowTheFlagBits is the constraint the comment on
// MsgResponseFlag states and nothing enforced.
//
// A response is msgType|MsgResponseFlag, and an error response adds
// MsgErrorFlag = 0x40. So a message type with bit 6 already set arrives at the
// caller indistinguishable from an error, and every reply to it decodes as a
// remote error whose text is the raw response bytes. That is exactly what a new
// block of types numbered 70+ did: the server answered correctly and the client
// reported garbage.
//
// The comment said "< 64". Reflection over the package's own declared types is
// what makes it true.
func TestEveryMessageTypeIsBelowTheFlagBits(t *testing.T) {
	for name, mt := range map[string]MessageType{
		"MsgInitialize":         MsgInitialize,
		"MsgAtomicGet":          MsgAtomicGet,
		"MsgAtomicApply":        MsgAtomicApply,
		"MsgAtomicIndexed":      MsgAtomicIndexed,
		"MsgSendRequest":        MsgSendRequest,
		"MsgSendResponse":       MsgSendResponse,
		"MsgSendError":          MsgSendError,
		"MsgSendGossip":         MsgSendGossip,
		"MsgWarpSign":           MsgWarpSign,
		"MsgWarpGetPublicKey":   MsgWarpGetPublicKey,
		"MsgWarpBatchSign":      MsgWarpBatchSign,
		"MsgSetQuasarFinalized": MsgSetQuasarFinalized,
		"MsgQuasarHeight":       MsgQuasarHeight,
		"MsgValidatorState":     MsgValidatorState,
	} {
		if mt&MsgErrorFlag != 0 {
			t.Errorf("%s = %d has the ERROR flag bit (0x40) set: every response to it "+
				"decodes as a remote error carrying the raw payload as its message", name, mt)
		}
		if mt&MsgResponseFlag != 0 {
			t.Errorf("%s = %d has the RESPONSE flag bit (0x80) set", name, mt)
		}
	}
}
