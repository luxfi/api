// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

// SendRequestMsg contains request to send to nodes (p2p.Sender.SendRequest)
type SendRequestMsg struct {
	NodeIDs   [][]byte
	RequestID uint32
	Request   []byte // Zero-copy payload
}

// Encode serializes SendRequestMsg to the buffer
func (m *SendRequestMsg) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.NodeIDs)))
	for _, nodeID := range m.NodeIDs {
		buf.WriteBytes(nodeID)
	}
	buf.WriteUint32(m.RequestID)
	buf.WriteBytes(m.Request)
}

// Decode deserializes SendRequestMsg from the reader
func (m *SendRequestMsg) Decode(r *Reader) error {
	count, err := r.ReadUint32()
	if err != nil {
		return err
	}
	m.NodeIDs = make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		if m.NodeIDs[i], err = r.ReadBytes(); err != nil {
			return err
		}
	}
	if m.RequestID, err = r.ReadUint32(); err != nil {
		return err
	}
	m.Request, err = r.ReadBytes()
	return err
}

// SendResponseMsg contains response to send to a node (p2p.Sender.SendResponse)
type SendResponseMsg struct {
	NodeID    []byte
	RequestID uint32
	Response  []byte // Zero-copy payload
}

// Encode serializes SendResponseMsg to the buffer
func (m *SendResponseMsg) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteUint32(m.RequestID)
	buf.WriteBytes(m.Response)
}

// Decode deserializes SendResponseMsg from the reader
func (m *SendResponseMsg) Decode(r *Reader) error {
	var err error
	if m.NodeID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.RequestID, err = r.ReadUint32(); err != nil {
		return err
	}
	m.Response, err = r.ReadBytes()
	return err
}

// SendErrorMsg contains error to send to a node (p2p.Sender.SendError)
type SendErrorMsg struct {
	NodeID       []byte
	RequestID    uint32
	ErrorCode    int32
	ErrorMessage string
}

// Encode serializes SendErrorMsg to the buffer
func (m *SendErrorMsg) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteUint32(m.RequestID)
	buf.WriteInt32(m.ErrorCode)
	buf.WriteString(m.ErrorMessage)
}

// Decode deserializes SendErrorMsg from the reader
func (m *SendErrorMsg) Decode(r *Reader) error {
	var err error
	if m.NodeID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.RequestID, err = r.ReadUint32(); err != nil {
		return err
	}
	if m.ErrorCode, err = r.ReadInt32(); err != nil {
		return err
	}
	m.ErrorMessage, err = r.ReadString()
	return err
}

// SendGossipMsg contains gossip message to send (p2p.Sender.SendGossip)
type SendGossipMsg struct {
	NodeIDs       [][]byte
	Validators    uint64
	NonValidators uint64
	Peers         uint64
	Msg           []byte // Zero-copy gossip payload
}

// Encode serializes SendGossipMsg to the buffer
func (m *SendGossipMsg) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.NodeIDs)))
	for _, nodeID := range m.NodeIDs {
		buf.WriteBytes(nodeID)
	}
	buf.WriteUint64(m.Validators)
	buf.WriteUint64(m.NonValidators)
	buf.WriteUint64(m.Peers)
	buf.WriteBytes(m.Msg)
}

// Decode deserializes SendGossipMsg from the reader
func (m *SendGossipMsg) Decode(r *Reader) error {
	count, err := r.ReadUint32()
	if err != nil {
		return err
	}
	m.NodeIDs = make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		if m.NodeIDs[i], err = r.ReadBytes(); err != nil {
			return err
		}
	}
	if m.Validators, err = r.ReadUint64(); err != nil {
		return err
	}
	if m.NonValidators, err = r.ReadUint64(); err != nil {
		return err
	}
	if m.Peers, err = r.ReadUint64(); err != nil {
		return err
	}
	m.Msg, err = r.ReadBytes()
	return err
}
