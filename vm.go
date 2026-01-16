// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/luxfi/consensus/runtime"
	"github.com/luxfi/ids"
)

// VM is the minimal interface that all VMs must implement.
// This is the execution boundary between consensus and VM logic.
type VM interface {
	// Initialize is called when the VM is first created.
	// rt contains chain wiring (IDs, logger, db, etc.)
	// ctx is used only for this initialization call's cancellation.
	Initialize(
		ctx context.Context,
		rt *runtime.Runtime,
		dbManager interface{}, // database.Manager
		genesisBytes []byte,
		upgradeBytes []byte,
		configBytes []byte,
		toEngine chan<- Message,
		fxs []*Fx,
		appSender AppSender,
	) error

	// SetState is called to notify the VM of its execution state
	SetState(ctx context.Context, state State) error

	// Shutdown is called when the VM should stop
	Shutdown(ctx context.Context) error

	// Version returns the VM's version string
	Version(ctx context.Context) (string, error)

	// CreateHandlers returns HTTP handlers for this VM's API
	CreateHandlers(ctx context.Context) (map[string]http.Handler, error)
}

// ChainVM is an optional interface for block-based VMs
type ChainVM interface {
	VM

	// BuildBlock builds a new block
	BuildBlock(ctx context.Context) (Block, error)

	// ParseBlock parses a block from bytes
	ParseBlock(ctx context.Context, blockBytes []byte) (Block, error)

	// GetBlock returns a block by ID
	GetBlock(ctx context.Context, blkID ids.ID) (Block, error)

	// SetPreference sets the preferred block
	SetPreference(ctx context.Context, blkID ids.ID) error

	// LastAccepted returns the last accepted block ID
	LastAccepted(ctx context.Context) (ids.ID, error)
}

// Block represents a block in the blockchain
type Block interface {
	// ID returns the block's unique identifier
	ID() ids.ID

	// Parent returns the parent block's ID
	Parent() ids.ID

	// Height returns the block's height
	Height() uint64

	// Timestamp returns the block's timestamp
	Timestamp() time.Time

	// Bytes returns the block's serialized form
	Bytes() []byte

	// Verify verifies the block is valid
	Verify(ctx context.Context) error

	// Accept marks the block as accepted
	Accept(ctx context.Context) error

	// Reject marks the block as rejected
	Reject(ctx context.Context) error

	// Status returns the block's status
	Status() Status
}

// Status represents a block's status
type Status uint32

const (
	StatusUnknown Status = iota
	StatusProcessing
	StatusRejected
	StatusAccepted
)

// State represents the VM's execution state
type State uint8

const (
	StateBootstrapping State = iota
	StateNormalOp
)

// Message is sent from VM to consensus engine
type Message struct {
	Type MessageType
}

// MessageType identifies the type of message
type MessageType uint8

const (
	MessagePendingTxs MessageType = iota
	MessageStateSyncDone
)

// Fx represents a feature extension
type Fx struct {
	ID ids.ID
	Fx interface{}
}

// AppSender sends application-level messages
type AppSender interface {
	SendAppRequest(ctx context.Context, nodeIDs []ids.NodeID, requestID uint32, request []byte) error
	SendAppResponse(ctx context.Context, nodeID ids.NodeID, requestID uint32, response []byte) error
	SendAppError(ctx context.Context, nodeID ids.NodeID, requestID uint32, errorCode int32, errorMessage string) error
	SendAppGossip(ctx context.Context, config SendConfig, msg []byte) error
	SendAppGossipSpecific(ctx context.Context, nodeIDs []ids.NodeID, msg []byte) error
	SendCrossChainAppRequest(ctx context.Context, chainID ids.ID, requestID uint32, msg []byte) error
	SendCrossChainAppResponse(ctx context.Context, chainID ids.ID, requestID uint32, msg []byte) error
	SendCrossChainAppError(ctx context.Context, chainID ids.ID, requestID uint32, errorCode int32, errorMessage string) error
}

// SendConfig configures gossip sending behavior
type SendConfig struct {
	Validators    int
	NonValidators int
	Peers         int
}
