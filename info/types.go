// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package info

import (
	"net/netip"

	"github.com/luxfi/api/types"
	"github.com/luxfi/ids"
	"github.com/luxfi/math/set"
	"github.com/luxfi/p2p/peer"
)

// ProofOfPossession is a JSON-friendly representation of a BLS PoP.
type ProofOfPossession struct {
	PublicKey         string `json:"publicKey"`
	ProofOfPossession string `json:"proofOfPossession"`
}

// GetNodeVersionReply are the results from calling GetNodeVersion.
type GetNodeVersionReply struct {
	Version            string            `json:"version"`
	DatabaseVersion    string            `json:"databaseVersion"`
	RPCProtocolVersion types.Uint32      `json:"rpcProtocolVersion"`
	GitCommit          string            `json:"gitCommit"`
	VMVersions         map[string]string `json:"vmVersions"`
}

// GetNodeIDReply are the results from calling GetNodeID.
type GetNodeIDReply struct {
	NodeID  ids.NodeID         `json:"nodeID"`
	NodePOP *ProofOfPossession `json:"nodePOP"`
}

// GetNetworkIDReply are the results from calling GetNetworkID.
type GetNetworkIDReply struct {
	NetworkID types.Uint32 `json:"networkID"`
}

// GetNodeIPReply are the results from calling GetNodeIP.
type GetNodeIPReply struct {
	IP netip.AddrPort `json:"ip"`
}

// GetNetworkNameReply is the result from calling GetNetworkName.
type GetNetworkNameReply struct {
	NetworkName string `json:"networkName"`
}

// GetBlockchainIDArgs are the arguments for calling GetBlockchainID.
type GetBlockchainIDArgs struct {
	Alias string `json:"alias"`
}

// GetBlockchainIDReply are the results from calling GetBlockchainID.
type GetBlockchainIDReply struct {
	BlockchainID ids.ID `json:"blockchainID"`
}

// PeersArgs are the arguments for calling Peers.
type PeersArgs struct {
	NodeIDs []ids.NodeID `json:"nodeIDs"`
}

// Peer is information about a peer in the network.
type Peer struct {
	peer.Info

	Benched []string `json:"benched"`
}

// PeersReply are the results from calling Peers.
type PeersReply struct {
	NumPeers types.Uint64 `json:"numPeers"`
	Peers    []Peer       `json:"peers"`
}

// IsBootstrappedArgs are the arguments for calling IsBootstrapped.
type IsBootstrappedArgs struct {
	Chain string `json:"chain"`
}

// IsBootstrappedResponse are the results from calling IsBootstrapped.
type IsBootstrappedResponse struct {
	IsBootstrapped bool `json:"isBootstrapped"`
}

// UptimeResponse are the results from calling Uptime.
type UptimeResponse struct {
	RewardingStakePercentage  types.Float64 `json:"rewardingStakePercentage"`
	WeightedAveragePercentage types.Float64 `json:"weightedAveragePercentage"`
}

// LP is information about an LP proposal.
type LP struct {
	SupportWeight types.Uint64        `json:"supportWeight"`
	Supporters    set.Set[ids.NodeID] `json:"supporters"`
	ObjectWeight  types.Uint64        `json:"objectWeight"`
	Objectors     set.Set[ids.NodeID] `json:"objectors"`
	AbstainWeight types.Uint64        `json:"abstainWeight"`
}

// LPsReply are the results from calling LPs.
type LPsReply struct {
	LPs map[uint32]*LP `json:"lps"`
}

// GetTxFeeResponse are the results from calling GetTxFee.
type GetTxFeeResponse struct {
	TxFee                  types.Uint64 `json:"txFee"`
	CreateAssetTxFee       types.Uint64 `json:"createAssetTxFee"`
	CreateNetworkTxFee     types.Uint64 `json:"createNetworkTxFee"`
	TransformChainTxFee    types.Uint64 `json:"transformChainTxFee"`
	CreateChainTxFee       types.Uint64 `json:"createChainTxFee"`
	AddNetworkValidatorFee types.Uint64 `json:"addNetworkValidatorFee"`
	AddNetworkDelegatorFee types.Uint64 `json:"addNetworkDelegatorFee"`
}

// GetVMsReply contains the response metadata for GetVMs.
type GetVMsReply struct {
	VMs map[ids.ID][]string `json:"vms"`
	Fxs map[ids.ID]string   `json:"fxs"`
}
