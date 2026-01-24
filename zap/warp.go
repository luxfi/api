// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

// WarpSignRequest encodes a warp signing request
type WarpSignRequest struct {
	NetworkID     uint32
	SourceChainID []byte // 32 bytes
	Payload       []byte
}

// Encode writes the request to the buffer
func (r *WarpSignRequest) Encode(buf *Buffer) {
	buf.WriteUint32(r.NetworkID)
	buf.WriteBytes(r.SourceChainID)
	buf.WriteBytes(r.Payload)
}

// Decode reads the request from the reader
func (r *WarpSignRequest) Decode(rd *Reader) error {
	var err error
	r.NetworkID, err = rd.ReadUint32()
	if err != nil {
		return err
	}
	r.SourceChainID, err = rd.ReadBytes()
	if err != nil {
		return err
	}
	r.Payload, err = rd.ReadBytes()
	return err
}

// WarpSignResponse encodes a warp signing response
type WarpSignResponse struct {
	Signature []byte
	Error     string
}

// Encode writes the response to the buffer
func (r *WarpSignResponse) Encode(buf *Buffer) {
	buf.WriteBytes(r.Signature)
	buf.WriteString(r.Error)
}

// Decode reads the response from the reader
func (r *WarpSignResponse) Decode(rd *Reader) error {
	var err error
	r.Signature, err = rd.ReadBytes()
	if err != nil {
		return err
	}
	r.Error, err = rd.ReadString()
	return err
}

// WarpGetPublicKeyRequest encodes a get public key request
type WarpGetPublicKeyRequest struct{}

// Encode writes the request to the buffer
func (r *WarpGetPublicKeyRequest) Encode(buf *Buffer) {
	// No fields
}

// Decode reads the request from the reader
func (r *WarpGetPublicKeyRequest) Decode(rd *Reader) error {
	return nil
}

// WarpGetPublicKeyResponse encodes a get public key response
type WarpGetPublicKeyResponse struct {
	PublicKey []byte
	Error     string
}

// Encode writes the response to the buffer
func (r *WarpGetPublicKeyResponse) Encode(buf *Buffer) {
	buf.WriteBytes(r.PublicKey)
	buf.WriteString(r.Error)
}

// Decode reads the response from the reader
func (r *WarpGetPublicKeyResponse) Decode(rd *Reader) error {
	var err error
	r.PublicKey, err = rd.ReadBytes()
	if err != nil {
		return err
	}
	r.Error, err = rd.ReadString()
	return err
}

// WarpBatchSignRequest encodes a batch signing request
type WarpBatchSignRequest struct {
	Messages []WarpSignRequest
}

// Encode writes the request to the buffer
func (r *WarpBatchSignRequest) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(r.Messages)))
	for i := range r.Messages {
		r.Messages[i].Encode(buf)
	}
}

// Decode reads the request from the reader
func (r *WarpBatchSignRequest) Decode(rd *Reader) error {
	count, err := rd.ReadUint32()
	if err != nil {
		return err
	}
	r.Messages = make([]WarpSignRequest, count)
	for i := uint32(0); i < count; i++ {
		if err := r.Messages[i].Decode(rd); err != nil {
			return err
		}
	}
	return nil
}

// WarpBatchSignResponse encodes a batch signing response
type WarpBatchSignResponse struct {
	Signatures [][]byte
	Errors     []string
}

// Encode writes the response to the buffer
func (r *WarpBatchSignResponse) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(r.Signatures)))
	for _, sig := range r.Signatures {
		buf.WriteBytes(sig)
	}
	buf.WriteUint32(uint32(len(r.Errors)))
	for _, errStr := range r.Errors {
		buf.WriteString(errStr)
	}
}

// Decode reads the response from the reader
func (r *WarpBatchSignResponse) Decode(rd *Reader) error {
	sigCount, err := rd.ReadUint32()
	if err != nil {
		return err
	}
	r.Signatures = make([][]byte, sigCount)
	for i := uint32(0); i < sigCount; i++ {
		r.Signatures[i], err = rd.ReadBytes()
		if err != nil {
			return err
		}
	}
	errCount, err := rd.ReadUint32()
	if err != nil {
		return err
	}
	r.Errors = make([]string, errCount)
	for i := uint32(0); i < errCount; i++ {
		r.Errors[i], err = rd.ReadString()
		if err != nil {
			return err
		}
	}
	return nil
}
