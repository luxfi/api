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
	// Consensus describes the active consensus configuration. Populated
	// from the in-memory consensus engine state at request time. Omitted
	// for legacy clients that don't wire it through.
	Consensus *ConsensusInfo `json:"consensus,omitempty"`
}

// ConsensusInfo summarises the consensus configuration of a running node so
// callers don't have to scrape the boot logs to learn whether Quasar is in
// triple/dual/classical mode.
type ConsensusInfo struct {
	// Mode is one of "triple" (BLS + Ringtail + ML-DSA), "dual" (BLS +
	// Ringtail), or "classical" (BLS only). Free-form so future modes
	// don't bump the API version.
	Mode string `json:"mode"`
	// BLS is always true for production Quasar nodes.
	BLS bool `json:"bls"`
	// Ringtail is true when the post-quantum lattice threshold path is wired.
	Ringtail bool `json:"ringtail"`
	// MLDSA is true when ML-DSA-65 (FIPS 204) signature verification is wired.
	MLDSA bool `json:"mlDSA"`
	// PlatformVM is true when the production PlatformVM wiring is in use,
	// false when running with the dev/stub PlatformVM.
	PlatformVM bool `json:"platformVM"`
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
