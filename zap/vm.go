// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zap

// State represents the VM state
type State uint8

const (
	StateUnspecified State = iota
	StateStateSyncing
	StateBootstrapping
	StateNormalOp
)

// Error represents VM errors
type Error uint8

const (
	ErrorUnspecified Error = iota
	ErrorClosed
	ErrorNotFound
	ErrorStateSyncNotImplemented
)

// NetworkUpgrades contains network upgrade timestamps
type NetworkUpgrades struct {
	ApricotPhase1Time            int64
	ApricotPhase2Time            int64
	ApricotPhase3Time            int64
	ApricotPhase4Time            int64
	ApricotPhase4MinPChainHeight uint64
	ApricotPhase5Time            int64
	ApricotPhasePre6Time         int64
	ApricotPhase6Time            int64
	ApricotPhasePost6Time        int64
	BanffTime                    int64
	CortinaTime                  int64
	CortinaXChainStopVertexID    []byte
	DurangoTime                  int64
	EtnaTime                     int64
	FortunaTime                  int64
	GraniteTime                  int64
}

// Encode serializes NetworkUpgrades to the buffer
func (n *NetworkUpgrades) Encode(buf *Buffer) {
	buf.WriteInt64(n.ApricotPhase1Time)
	buf.WriteInt64(n.ApricotPhase2Time)
	buf.WriteInt64(n.ApricotPhase3Time)
	buf.WriteInt64(n.ApricotPhase4Time)
	buf.WriteUint64(n.ApricotPhase4MinPChainHeight)
	buf.WriteInt64(n.ApricotPhase5Time)
	buf.WriteInt64(n.ApricotPhasePre6Time)
	buf.WriteInt64(n.ApricotPhase6Time)
	buf.WriteInt64(n.ApricotPhasePost6Time)
	buf.WriteInt64(n.BanffTime)
	buf.WriteInt64(n.CortinaTime)
	buf.WriteBytes(n.CortinaXChainStopVertexID)
	buf.WriteInt64(n.DurangoTime)
	buf.WriteInt64(n.EtnaTime)
	buf.WriteInt64(n.FortunaTime)
	buf.WriteInt64(n.GraniteTime)
}

// Decode deserializes NetworkUpgrades from the reader
func (n *NetworkUpgrades) Decode(r *Reader) error {
	var err error
	if n.ApricotPhase1Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhase2Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhase3Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhase4Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhase4MinPChainHeight, err = r.ReadUint64(); err != nil {
		return err
	}
	if n.ApricotPhase5Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhasePre6Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhase6Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.ApricotPhasePost6Time, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.BanffTime, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.CortinaTime, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.CortinaXChainStopVertexID, err = r.ReadBytes(); err != nil {
		return err
	}
	if n.DurangoTime, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.EtnaTime, err = r.ReadInt64(); err != nil {
		return err
	}
	if n.FortunaTime, err = r.ReadInt64(); err != nil {
		return err
	}
	n.GraniteTime, err = r.ReadInt64()
	return err
}

// InitializeRequest contains initialization parameters
type InitializeRequest struct {
	NetworkID       uint32
	ChainID         []byte
	NodeID          []byte
	PublicKey       []byte
	XChainID        []byte
	CChainID        []byte
	LuxAssetID      []byte
	ChainDataDir    string
	GenesisBytes    []byte
	UpgradeBytes    []byte
	ConfigBytes     []byte
	DBServerAddr    string
	ServerAddr      string
	NetworkUpgrades NetworkUpgrades
}

// Encode serializes InitializeRequest to the buffer
func (m *InitializeRequest) Encode(buf *Buffer) {
	buf.WriteUint32(m.NetworkID)
	buf.WriteBytes(m.ChainID)
	buf.WriteBytes(m.NodeID)
	buf.WriteBytes(m.PublicKey)
	buf.WriteBytes(m.XChainID)
	buf.WriteBytes(m.CChainID)
	buf.WriteBytes(m.LuxAssetID)
	buf.WriteString(m.ChainDataDir)
	buf.WriteBytes(m.GenesisBytes)
	buf.WriteBytes(m.UpgradeBytes)
	buf.WriteBytes(m.ConfigBytes)
	buf.WriteString(m.DBServerAddr)
	buf.WriteString(m.ServerAddr)
	m.NetworkUpgrades.Encode(buf)
}

// Decode deserializes InitializeRequest from the reader
func (m *InitializeRequest) Decode(r *Reader) error {
	var err error
	if m.NetworkID, err = r.ReadUint32(); err != nil {
		return err
	}
	if m.ChainID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.NodeID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.PublicKey, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.XChainID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.CChainID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.LuxAssetID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.ChainDataDir, err = r.ReadString(); err != nil {
		return err
	}
	if m.GenesisBytes, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.UpgradeBytes, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.ConfigBytes, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.DBServerAddr, err = r.ReadString(); err != nil {
		return err
	}
	if m.ServerAddr, err = r.ReadString(); err != nil {
		return err
	}
	return m.NetworkUpgrades.Decode(r)
}

// InitializeResponse contains initialization results
type InitializeResponse struct {
	LastAcceptedID       []byte
	LastAcceptedParentID []byte
	Height               uint64
	Bytes                []byte
	Timestamp            int64
}

// Encode serializes InitializeResponse to the buffer
func (m *InitializeResponse) Encode(buf *Buffer) {
	buf.WriteBytes(m.LastAcceptedID)
	buf.WriteBytes(m.LastAcceptedParentID)
	buf.WriteUint64(m.Height)
	buf.WriteBytes(m.Bytes)
	buf.WriteInt64(m.Timestamp)
}

// Decode deserializes InitializeResponse from the reader
func (m *InitializeResponse) Decode(r *Reader) error {
	var err error
	if m.LastAcceptedID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.LastAcceptedParentID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.Height, err = r.ReadUint64(); err != nil {
		return err
	}
	if m.Bytes, err = r.ReadBytes(); err != nil {
		return err
	}
	m.Timestamp, err = r.ReadInt64()
	return err
}

// SetStateRequest contains state change request
type SetStateRequest struct {
	State State
}

// Encode serializes SetStateRequest to the buffer
func (m *SetStateRequest) Encode(buf *Buffer) {
	buf.WriteUint8(uint8(m.State))
}

// Decode deserializes SetStateRequest from the reader
func (m *SetStateRequest) Decode(r *Reader) error {
	v, err := r.ReadUint8()
	m.State = State(v)
	return err
}

// SetStateResponse contains state change results
type SetStateResponse struct {
	LastAcceptedID       []byte
	LastAcceptedParentID []byte
	Height               uint64
	Bytes                []byte
	Timestamp            int64
}

// Encode serializes SetStateResponse to the buffer
func (m *SetStateResponse) Encode(buf *Buffer) {
	buf.WriteBytes(m.LastAcceptedID)
	buf.WriteBytes(m.LastAcceptedParentID)
	buf.WriteUint64(m.Height)
	buf.WriteBytes(m.Bytes)
	buf.WriteInt64(m.Timestamp)
}

// Decode deserializes SetStateResponse from the reader
func (m *SetStateResponse) Decode(r *Reader) error {
	var err error
	if m.LastAcceptedID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.LastAcceptedParentID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.Height, err = r.ReadUint64(); err != nil {
		return err
	}
	if m.Bytes, err = r.ReadBytes(); err != nil {
		return err
	}
	m.Timestamp, err = r.ReadInt64()
	return err
}

// BuildBlockRequest contains block building parameters
type BuildBlockRequest struct {
	PChainHeight    uint64
	HasPChainHeight bool
}

// Encode serializes BuildBlockRequest to the buffer
func (m *BuildBlockRequest) Encode(buf *Buffer) {
	buf.WriteBool(m.HasPChainHeight)
	if m.HasPChainHeight {
		buf.WriteUint64(m.PChainHeight)
	}
}

// Decode deserializes BuildBlockRequest from the reader
func (m *BuildBlockRequest) Decode(r *Reader) error {
	var err error
	if m.HasPChainHeight, err = r.ReadBool(); err != nil {
		return err
	}
	if m.HasPChainHeight {
		m.PChainHeight, err = r.ReadUint64()
	}
	return err
}

// BlockResponse contains block data (used by BuildBlock, ParseBlock, GetBlock)
type BlockResponse struct {
	ID                []byte
	ParentID          []byte
	Bytes             []byte // Zero-copy block data
	Height            uint64
	Timestamp         int64
	VerifyWithContext bool
	Err               Error
}

// Encode serializes BlockResponse to the buffer
func (m *BlockResponse) Encode(buf *Buffer) {
	buf.WriteBytes(m.ID)
	buf.WriteBytes(m.ParentID)
	buf.WriteBytes(m.Bytes)
	buf.WriteUint64(m.Height)
	buf.WriteInt64(m.Timestamp)
	buf.WriteBool(m.VerifyWithContext)
	buf.WriteUint8(uint8(m.Err))
}

// Decode deserializes BlockResponse from the reader
func (m *BlockResponse) Decode(r *Reader) error {
	var err error
	if m.ID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.ParentID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.Bytes, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.Height, err = r.ReadUint64(); err != nil {
		return err
	}
	if m.Timestamp, err = r.ReadInt64(); err != nil {
		return err
	}
	if m.VerifyWithContext, err = r.ReadBool(); err != nil {
		return err
	}
	v, err := r.ReadUint8()
	m.Err = Error(v)
	return err
}

// ParseBlockRequest contains bytes to parse
type ParseBlockRequest struct {
	Bytes []byte // Zero-copy input
}

// Encode serializes ParseBlockRequest to the buffer
func (m *ParseBlockRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.Bytes)
}

// Decode deserializes ParseBlockRequest from the reader
func (m *ParseBlockRequest) Decode(r *Reader) error {
	var err error
	m.Bytes, err = r.ReadBytes()
	return err
}

// GetBlockRequest contains block ID to retrieve
type GetBlockRequest struct {
	ID []byte
}

// Encode serializes GetBlockRequest to the buffer
func (m *GetBlockRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.ID)
}

// Decode deserializes GetBlockRequest from the reader
func (m *GetBlockRequest) Decode(r *Reader) error {
	var err error
	m.ID, err = r.ReadBytes()
	return err
}

// SetPreferenceRequest contains preferred block ID
type SetPreferenceRequest struct {
	ID []byte
}

// Encode serializes SetPreferenceRequest to the buffer
func (m *SetPreferenceRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.ID)
}

// Decode deserializes SetPreferenceRequest from the reader
func (m *SetPreferenceRequest) Decode(r *Reader) error {
	var err error
	m.ID, err = r.ReadBytes()
	return err
}

// BlockVerifyRequest contains block verification parameters
type BlockVerifyRequest struct {
	Bytes           []byte
	PChainHeight    uint64
	HasPChainHeight bool
}

// Encode serializes BlockVerifyRequest to the buffer
func (m *BlockVerifyRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.Bytes)
	buf.WriteBool(m.HasPChainHeight)
	if m.HasPChainHeight {
		buf.WriteUint64(m.PChainHeight)
	}
}

// Decode deserializes BlockVerifyRequest from the reader
func (m *BlockVerifyRequest) Decode(r *Reader) error {
	var err error
	if m.Bytes, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.HasPChainHeight, err = r.ReadBool(); err != nil {
		return err
	}
	if m.HasPChainHeight {
		m.PChainHeight, err = r.ReadUint64()
	}
	return err
}

// BlockVerifyResponse contains verification result
type BlockVerifyResponse struct {
	Timestamp int64
}

// Encode serializes BlockVerifyResponse to the buffer
func (m *BlockVerifyResponse) Encode(buf *Buffer) {
	buf.WriteInt64(m.Timestamp)
}

// Decode deserializes BlockVerifyResponse from the reader
func (m *BlockVerifyResponse) Decode(r *Reader) error {
	var err error
	m.Timestamp, err = r.ReadInt64()
	return err
}

// BlockAcceptRequest contains block ID to accept
type BlockAcceptRequest struct {
	ID []byte
}

// Encode serializes BlockAcceptRequest to the buffer
func (m *BlockAcceptRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.ID)
}

// Decode deserializes BlockAcceptRequest from the reader
func (m *BlockAcceptRequest) Decode(r *Reader) error {
	var err error
	m.ID, err = r.ReadBytes()
	return err
}

// BlockRejectRequest contains block ID to reject
type BlockRejectRequest struct {
	ID []byte
}

// Encode serializes BlockRejectRequest to the buffer
func (m *BlockRejectRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.ID)
}

// Decode deserializes BlockRejectRequest from the reader
func (m *BlockRejectRequest) Decode(r *Reader) error {
	var err error
	m.ID, err = r.ReadBytes()
	return err
}

// BatchedParseBlockRequest contains multiple blocks to parse
type BatchedParseBlockRequest struct {
	Requests [][]byte
}

// Encode serializes BatchedParseBlockRequest to the buffer
func (m *BatchedParseBlockRequest) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.Requests)))
	for _, req := range m.Requests {
		buf.WriteBytes(req)
	}
}

// Decode deserializes BatchedParseBlockRequest from the reader
func (m *BatchedParseBlockRequest) Decode(r *Reader) error {
	count, err := r.ReadUint32()
	if err != nil {
		return err
	}
	m.Requests = make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		if m.Requests[i], err = r.ReadBytes(); err != nil {
			return err
		}
	}
	return nil
}

// BatchedParseBlockResponse contains parsed blocks
type BatchedParseBlockResponse struct {
	Responses []BlockResponse
}

// Encode serializes BatchedParseBlockResponse to the buffer
func (m *BatchedParseBlockResponse) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.Responses)))
	for i := range m.Responses {
		m.Responses[i].Encode(buf)
	}
}

// Decode deserializes BatchedParseBlockResponse from the reader
func (m *BatchedParseBlockResponse) Decode(r *Reader) error {
	count, err := r.ReadUint32()
	if err != nil {
		return err
	}
	m.Responses = make([]BlockResponse, count)
	for i := uint32(0); i < count; i++ {
		if err := m.Responses[i].Decode(r); err != nil {
			return err
		}
	}
	return nil
}

// GetAncestorsRequest contains ancestor retrieval parameters
type GetAncestorsRequest struct {
	BlkID                 []byte
	MaxBlocksNum          int32
	MaxBlocksSize         int32
	MaxBlocksRetrivalTime int64
}

// Encode serializes GetAncestorsRequest to the buffer
func (m *GetAncestorsRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.BlkID)
	buf.WriteInt32(m.MaxBlocksNum)
	buf.WriteInt32(m.MaxBlocksSize)
	buf.WriteInt64(m.MaxBlocksRetrivalTime)
}

// Decode deserializes GetAncestorsRequest from the reader
func (m *GetAncestorsRequest) Decode(r *Reader) error {
	var err error
	if m.BlkID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.MaxBlocksNum, err = r.ReadInt32(); err != nil {
		return err
	}
	if m.MaxBlocksSize, err = r.ReadInt32(); err != nil {
		return err
	}
	m.MaxBlocksRetrivalTime, err = r.ReadInt64()
	return err
}

// GetAncestorsResponse contains ancestor blocks
type GetAncestorsResponse struct {
	BlksBytes [][]byte
}

// Encode serializes GetAncestorsResponse to the buffer
func (m *GetAncestorsResponse) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.BlksBytes)))
	for _, blk := range m.BlksBytes {
		buf.WriteBytes(blk)
	}
}

// Decode deserializes GetAncestorsResponse from the reader
func (m *GetAncestorsResponse) Decode(r *Reader) error {
	count, err := r.ReadUint32()
	if err != nil {
		return err
	}
	m.BlksBytes = make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		if m.BlksBytes[i], err = r.ReadBytes(); err != nil {
			return err
		}
	}
	return nil
}

// GetBlockIDAtHeightRequest contains height to query
type GetBlockIDAtHeightRequest struct {
	Height uint64
}

// Encode serializes GetBlockIDAtHeightRequest to the buffer
func (m *GetBlockIDAtHeightRequest) Encode(buf *Buffer) {
	buf.WriteUint64(m.Height)
}

// Decode deserializes GetBlockIDAtHeightRequest from the reader
func (m *GetBlockIDAtHeightRequest) Decode(r *Reader) error {
	var err error
	m.Height, err = r.ReadUint64()
	return err
}

// GetBlockIDAtHeightResponse contains block ID at height
type GetBlockIDAtHeightResponse struct {
	BlkID []byte
	Err   Error
}

// Encode serializes GetBlockIDAtHeightResponse to the buffer
func (m *GetBlockIDAtHeightResponse) Encode(buf *Buffer) {
	buf.WriteBytes(m.BlkID)
	buf.WriteUint8(uint8(m.Err))
}

// Decode deserializes GetBlockIDAtHeightResponse from the reader
func (m *GetBlockIDAtHeightResponse) Decode(r *Reader) error {
	var err error
	if m.BlkID, err = r.ReadBytes(); err != nil {
		return err
	}
	v, err := r.ReadUint8()
	m.Err = Error(v)
	return err
}

// RequestMsg contains incoming request data
type RequestMsg struct {
	NodeID    []byte
	RequestID uint32
	Deadline  int64
	Request   []byte
}

// Encode serializes RequestMsg to the buffer
func (m *RequestMsg) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteUint32(m.RequestID)
	buf.WriteInt64(m.Deadline)
	buf.WriteBytes(m.Request)
}

// Decode deserializes RequestMsg from the reader
func (m *RequestMsg) Decode(r *Reader) error {
	var err error
	if m.NodeID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.RequestID, err = r.ReadUint32(); err != nil {
		return err
	}
	if m.Deadline, err = r.ReadInt64(); err != nil {
		return err
	}
	m.Request, err = r.ReadBytes()
	return err
}

// RequestFailedMsg contains failed request info
type RequestFailedMsg struct {
	NodeID       []byte
	RequestID    uint32
	ErrorCode    int32
	ErrorMessage string
}

// Encode serializes RequestFailedMsg to the buffer
func (m *RequestFailedMsg) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteUint32(m.RequestID)
	buf.WriteInt32(m.ErrorCode)
	buf.WriteString(m.ErrorMessage)
}

// Decode deserializes RequestFailedMsg from the reader
func (m *RequestFailedMsg) Decode(r *Reader) error {
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

// ResponseMsg contains response data
type ResponseMsg struct {
	NodeID    []byte
	RequestID uint32
	Response  []byte
}

// Encode serializes ResponseMsg to the buffer
func (m *ResponseMsg) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteUint32(m.RequestID)
	buf.WriteBytes(m.Response)
}

// Decode deserializes ResponseMsg from the reader
func (m *ResponseMsg) Decode(r *Reader) error {
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

// GossipMsg contains gossip message
type GossipMsg struct {
	NodeID []byte
	Msg    []byte
}

// Encode serializes GossipMsg to the buffer
func (m *GossipMsg) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteBytes(m.Msg)
}

// Decode deserializes GossipMsg from the reader
func (m *GossipMsg) Decode(r *Reader) error {
	var err error
	if m.NodeID, err = r.ReadBytes(); err != nil {
		return err
	}
	m.Msg, err = r.ReadBytes()
	return err
}

// ConnectedRequest contains connection info
type ConnectedRequest struct {
	NodeID []byte
	Name   string
	Major  uint32
	Minor  uint32
	Patch  uint32
}

// Encode serializes ConnectedRequest to the buffer
func (m *ConnectedRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
	buf.WriteString(m.Name)
	buf.WriteUint32(m.Major)
	buf.WriteUint32(m.Minor)
	buf.WriteUint32(m.Patch)
}

// Decode deserializes ConnectedRequest from the reader
func (m *ConnectedRequest) Decode(r *Reader) error {
	var err error
	if m.NodeID, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.Name, err = r.ReadString(); err != nil {
		return err
	}
	if m.Major, err = r.ReadUint32(); err != nil {
		return err
	}
	if m.Minor, err = r.ReadUint32(); err != nil {
		return err
	}
	m.Patch, err = r.ReadUint32()
	return err
}

// DisconnectedRequest contains disconnection info
type DisconnectedRequest struct {
	NodeID []byte
}

// Encode serializes DisconnectedRequest to the buffer
func (m *DisconnectedRequest) Encode(buf *Buffer) {
	buf.WriteBytes(m.NodeID)
}

// Decode deserializes DisconnectedRequest from the reader
func (m *DisconnectedRequest) Decode(r *Reader) error {
	var err error
	m.NodeID, err = r.ReadBytes()
	return err
}

// HealthResponse contains health check result
type HealthResponse struct {
	Details []byte
}

// Encode serializes HealthResponse to the buffer
func (m *HealthResponse) Encode(buf *Buffer) {
	buf.WriteBytes(m.Details)
}

// Decode deserializes HealthResponse from the reader
func (m *HealthResponse) Decode(r *Reader) error {
	var err error
	m.Details, err = r.ReadBytes()
	return err
}

// VersionResponse contains version info
type VersionResponse struct {
	Version string
}

// Encode serializes VersionResponse to the buffer
func (m *VersionResponse) Encode(buf *Buffer) {
	buf.WriteString(m.Version)
}

// Decode deserializes VersionResponse from the reader
func (m *VersionResponse) Decode(r *Reader) error {
	var err error
	m.Version, err = r.ReadString()
	return err
}

// WaitForEventResponse contains event type
type WaitForEventResponse struct {
	Message uint8
}

// Encode serializes WaitForEventResponse to the buffer
func (m *WaitForEventResponse) Encode(buf *Buffer) {
	buf.WriteUint8(m.Message)
}

// Decode deserializes WaitForEventResponse from the reader
func (m *WaitForEventResponse) Decode(r *Reader) error {
	var err error
	m.Message, err = r.ReadUint8()
	return err
}

// HTTPHandler contains HTTP handler info
type HTTPHandler struct {
	Prefix     string
	ServerAddr string
}

// Encode serializes HTTPHandler to the buffer
func (m *HTTPHandler) Encode(buf *Buffer) {
	buf.WriteString(m.Prefix)
	buf.WriteString(m.ServerAddr)
}

// Decode deserializes HTTPHandler from the reader
func (m *HTTPHandler) Decode(r *Reader) error {
	var err error
	if m.Prefix, err = r.ReadString(); err != nil {
		return err
	}
	m.ServerAddr, err = r.ReadString()
	return err
}

// CreateHandlersResponse contains handler list
type CreateHandlersResponse struct {
	Handlers []HTTPHandler
}

// Encode serializes CreateHandlersResponse to the buffer
func (m *CreateHandlersResponse) Encode(buf *Buffer) {
	buf.WriteUint32(uint32(len(m.Handlers)))
	for i := range m.Handlers {
		m.Handlers[i].Encode(buf)
	}
}

// Decode deserializes CreateHandlersResponse from the reader
func (m *CreateHandlersResponse) Decode(r *Reader) error {
	count, err := r.ReadUint32()
	if err != nil {
		return err
	}
	m.Handlers = make([]HTTPHandler, count)
	for i := uint32(0); i < count; i++ {
		if err := m.Handlers[i].Decode(r); err != nil {
			return err
		}
	}
	return nil
}

// NewHTTPHandlerResponse contains HTTP handler address
type NewHTTPHandlerResponse struct {
	ServerAddr string
}

// Encode serializes NewHTTPHandlerResponse to the buffer
func (m *NewHTTPHandlerResponse) Encode(buf *Buffer) {
	buf.WriteString(m.ServerAddr)
}

// Decode deserializes NewHTTPHandlerResponse from the reader
func (m *NewHTTPHandlerResponse) Decode(r *Reader) error {
	var err error
	m.ServerAddr, err = r.ReadString()
	return err
}
