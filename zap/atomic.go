// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

// atomic.go is the wire for the cross-chain atomic shared-memory surface — the
// primary-network import/export primitive platformvm, the dexvm and the DEX
// 0x9999 settlement seam all move value through.
//
// WHY IT IS ON THE WIRE AT ALL. One atomic.Memory exists per NODE; every chain
// gets a per-chain SharedMemory handle carved out of it. A VM linked into the
// node reaches its handle straight off runtime.Runtime. A VM hosted as a ZAP
// plugin — the C-Chain EVM and the dexvm both are — does not: the node rebuilds
// the plugin's Runtime from the InitializeRequest fields, and neither the handle
// (an interface over a live database) nor BCLookup can be copied into a struct
// literal. So the plugin's Runtime.SharedMemory was nil and every settlement
// reverted "requires the cross-chain atomic capability".
//
// DIRECTION AND FRAMING. The node binds a per-chain ZAP server over its own
// SharedMemory handle and names it in InitializeRequest.AtomicServerAddr; the
// plugin dials that addr as an ordinary ZAP client. The node is the SERVER, so
// the existing strictly request/response transport carries this unchanged — no
// duplex framing, no second demultiplexer, no gRPC. It is the identical shape
// the node already uses for DBServerAddr.
//
// TYPES. This package stays dependency-free on purpose (luxfi/vm imports
// luxfi/api, so the reverse would be a cycle). IDs are plain 32-byte slices and
// the request set is flattened; the plugin side converts to and from
// atomic.Requests.

// AtomicGetRequest asks the node for the values a peer chain published for this
// chain. It is SharedMemory.Get(peerChainID, keys) verbatim.
type AtomicGetRequest struct {
	PeerChainID []byte
	Keys        [][]byte
}

// Encode serializes AtomicGetRequest to the buffer
func (m *AtomicGetRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.PeerChainID)
	buf.WriteUint32(uint32(len(m.Keys)))
	for _, k := range m.Keys {
		buf.WriteBytes(k)
	}
}

// Decode deserializes AtomicGetRequest from the reader
func (m *AtomicGetRequest) Decode(r *Reader) error {
	var err error
	if m.PeerChainID, err = r.ReadBytes(); err != nil {
		return err
	}
	count, err := r.ReadCount()
	if err != nil {
		return err
	}
	m.Keys = make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		if m.Keys[i], err = r.ReadBytes(); err != nil {
			return err
		}
	}
	return nil
}

// AtomicGetResponse carries the values, positionally aligned with the request's
// keys. SharedMemory.Get guarantees len(values) == len(keys); an absent key is
// an EMPTY value at its index, never a short slice — the caller distinguishes
// "no such object" from "object present" by length, so that invariant is part of
// the wire contract and is asserted on decode by the caller.
type AtomicGetResponse struct {
	Values [][]byte
}

// Encode serializes AtomicGetResponse to the buffer
func (m *AtomicGetResponse) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.Values)))
	for _, v := range m.Values {
		buf.WriteBytes(v)
	}
}

// Decode deserializes AtomicGetResponse from the reader
func (m *AtomicGetResponse) Decode(r *Reader) error {
	count, err := r.ReadCount()
	if err != nil {
		return err
	}
	m.Values = make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		if m.Values[i], err = r.ReadBytes(); err != nil {
			return err
		}
	}
	return nil
}

// AtomicElement is one published cross-chain object: the value stored under Key,
// plus the traits it is indexed by. Mirrors atomic.Element field-for-field.
type AtomicElement struct {
	Key    []byte
	Value  []byte
	Traits [][]byte
}

// AtomicChainRequests is one peer chain's operations. Mirrors atomic.Requests.
type AtomicChainRequests struct {
	PeerChainID    []byte
	RemoveRequests [][]byte
	PutRequests    []*AtomicElement
}

// AtomicApplyRequest is SharedMemory.Apply(requests) — the atomic commit of a
// block's cross-chain puts and removes.
//
// NO BATCH CROSSES THIS WIRE. In-process, Apply takes the caller's database
// batch so the state commit and the shared-memory mutation land as one write; a
// batch is a live handle over a database the node owns and cannot be serialized.
// The plugin therefore relies on the property its flush window already has: the
// ops are derived from a monotone seq in CONSENSUS state, so a replayed window
// yields byte-identical Puts and Removes. Both are idempotent under replay, so
// at-least-once delivery of this message has exactly-once effect — which is why
// the caller must advance its flushed-seq marker only AFTER this call returns.
type AtomicApplyRequest struct {
	Requests []*AtomicChainRequests
}

// Encode serializes AtomicApplyRequest to the buffer
func (m *AtomicApplyRequest) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.Requests)))
	for _, req := range m.Requests {
		buf.WriteBytes(req.PeerChainID)
		buf.WriteUint32(uint32(len(req.RemoveRequests)))
		for _, rm := range req.RemoveRequests {
			buf.WriteBytes(rm)
		}
		buf.WriteUint32(uint32(len(req.PutRequests)))
		for _, e := range req.PutRequests {
			buf.WriteBytes(e.Key)
			buf.WriteBytes(e.Value)
			buf.WriteUint32(uint32(len(e.Traits)))
			for _, tr := range e.Traits {
				buf.WriteBytes(tr)
			}
		}
	}
}

// Decode deserializes AtomicApplyRequest from the reader
func (m *AtomicApplyRequest) Decode(r *Reader) error {
	chainCount, err := r.ReadCount()
	if err != nil {
		return err
	}
	m.Requests = make([]*AtomicChainRequests, chainCount)
	for i := uint32(0); i < chainCount; i++ {
		req := &AtomicChainRequests{}
		if req.PeerChainID, err = r.ReadBytes(); err != nil {
			return err
		}
		rmCount, err := r.ReadCount()
		if err != nil {
			return err
		}
		req.RemoveRequests = make([][]byte, rmCount)
		for j := uint32(0); j < rmCount; j++ {
			if req.RemoveRequests[j], err = r.ReadBytes(); err != nil {
				return err
			}
		}
		putCount, err := r.ReadCount()
		if err != nil {
			return err
		}
		req.PutRequests = make([]*AtomicElement, putCount)
		for j := uint32(0); j < putCount; j++ {
			e := &AtomicElement{}
			if e.Key, err = r.ReadBytes(); err != nil {
				return err
			}
			if e.Value, err = r.ReadBytes(); err != nil {
				return err
			}
			trCount, err := r.ReadCount()
			if err != nil {
				return err
			}
			e.Traits = make([][]byte, trCount)
			for k := uint32(0); k < trCount; k++ {
				if e.Traits[k], err = r.ReadBytes(); err != nil {
					return err
				}
			}
			req.PutRequests[j] = e
		}
		m.Requests[i] = req
	}
	return nil
}
