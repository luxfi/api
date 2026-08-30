// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// What a node will tell anyone who asks.
//
// SEVERAL SHAPES HERE ARE AN OBJECT ON THE JSON WIRE AND A LIST IN GO, and the
// reason is the same every time: a ZAP field is an offset, a map has no order to
// give one, so a reply carrying a map cannot cross between two processes at all.
// The list is ordered by the key it was an object of, which is an order a map
// never had — so the answer is now the same twice running as well.
//
// The JSON is unchanged. Each list marshals to the object it replaced, and reads
// one back.
package info

import (
	"encoding/json"
	"slices"

	"github.com/luxfi/api/types"
	"github.com/luxfi/ids"
)

// PeerInfo is the JSON wire shape for a peer in the Peers RPC reply. It is
// the local owner of this shape; previously this type came from
// github.com/luxfi/p2p/peer.Info, but that module is being archived per
// LP-201. Field tags are stable across the rename — JSON consumers see no
// difference.
type PeerInfo struct {
	// IP is the address this node reaches the peer at.
	IP types.Addr `json:"ip"`
	// PublicIP is the address the peer says it is reachable at.
	PublicIP types.Addr `json:"publicIP,omitempty"`
	// ID is the peer's node id.
	ID ids.NodeID `json:"nodeID"`
	// Version is the node software the peer is running.
	Version string `json:"version"`
	// LastSent is when this node last sent the peer a message.
	LastSent types.Time `json:"lastSent"`
	// LastReceived is when this node last heard from the peer.
	LastReceived types.Time `json:"lastReceived"`
	// ObservedUptime is the percentage of time this node has seen the peer up.
	ObservedUptime types.Uint32 `json:"observedUptime"`
	// TrackedChains are the chains the peer serves, ordered by id.
	TrackedChains []ids.ID `json:"trackedChains"`
	// SupportedLPs are the LPs the peer votes for, ascending.
	SupportedLPs []uint32 `json:"supportedLPs"`
	// ObjectedLPs are the LPs the peer votes against, ascending.
	ObjectedLPs []uint32 `json:"objectedLPs"`
}

// ProofOfPossession is a JSON-friendly representation of a BLS PoP.
type ProofOfPossession struct {
	PublicKey         string `json:"publicKey"`
	ProofOfPossession string `json:"proofOfPossession"`
}

// GetNodeVersionReply are the results from calling GetNodeVersion.
type GetNodeVersionReply struct {
	Version            string       `json:"version"`
	DatabaseVersion    string       `json:"databaseVersion"`
	RPCProtocolVersion types.Uint32 `json:"rpcProtocolVersion"`
	GitCommit          string       `json:"gitCommit"`
	VMVersions         VMVersions   `json:"vmVersions"`
	// Consensus describes the active consensus configuration. Populated
	// from the in-memory consensus engine state at request time. Omitted
	// for legacy clients that don't wire it through.
	Consensus *ConsensusInfo `json:"consensus,omitempty"`
}

// VMVersion is one VM and the version of it a node runs.
type VMVersion struct {
	// VM is the VM's alias.
	VM string `json:"vm"`
	// Version is the release of it this node runs.
	Version string `json:"version"`
}

// VMVersions is what a node runs of each VM: an object keyed by VM alias on the
// JSON wire, a list ordered by that alias here.
type VMVersions []VMVersion

func (v VMVersions) MarshalJSON() ([]byte, error) {
	m := make(map[string]string, len(v))
	for _, e := range v {
		m[e.VM] = e.Version
	}
	return json.Marshal(m)
}

func (v *VMVersions) UnmarshalJSON(b []byte) error {
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*v = make(VMVersions, 0, len(m))
	for vm, version := range m {
		*v = append(*v, VMVersion{VM: vm, Version: version})
	}
	slices.SortFunc(*v, func(x, y VMVersion) int { return cmpString(x.VM, y.VM) })
	return nil
}

// ConsensusInfo summarises the consensus configuration of a running node so
// callers don't have to scrape the boot logs to learn whether Quasar is in
// triple/dual/classical mode.
type ConsensusInfo struct {
	// Mode is one of "triple" (BLS + Corona + ML-DSA), "dual" (BLS +
	// Corona), or "classical" (BLS only). Free-form so future modes
	// don't bump the API version.
	Mode string `json:"mode"`
	// BLS is always true for production Quasar nodes.
	BLS bool `json:"bls"`
	// Corona is true when the post-quantum lattice threshold path is wired.
	Corona bool `json:"corona"`
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
	// IP is the address this node tells peers to reach it at.
	IP types.Addr `json:"ip"`
}

// GetNetworkNameReply is the result from calling GetNetworkName.
type GetNetworkNameReply struct {
	NetworkName string `json:"networkName"`
}

// GetBlockchainIDArgs are the arguments for calling GetBlockchainID.
type GetBlockchainIDArgs struct {
	// Alias is the name the chain answers to, such as "X".
	Alias string `json:"alias"`
}

// GetBlockchainIDReply are the results from calling GetBlockchainID.
type GetBlockchainIDReply struct {
	BlockchainID ids.ID `json:"blockchainID"`
}

// PeersArgs are the arguments for calling Peers.
type PeersArgs struct {
	// NodeIDs narrows the answer to these peers. Empty asks for all of them.
	NodeIDs []ids.NodeID `json:"nodeIDs"`
}

// Peer is information about a peer in the network.
type Peer struct {
	PeerInfo

	Benched []string `json:"benched"`
}

// PeersReply are the results from calling Peers.
type PeersReply struct {
	NumPeers types.Uint64 `json:"numPeers"`
	Peers    []Peer       `json:"peers"`
}

// IsBootstrappedArgs are the arguments for calling IsBootstrapped.
type IsBootstrappedArgs struct {
	// Chain is the alias or id of the chain to ask about.
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
	SupportWeight types.Uint64 `json:"supportWeight"`
	// Supporters are the peers voting for it, ordered by node id.
	Supporters   []ids.NodeID `json:"supporters"`
	ObjectWeight types.Uint64 `json:"objectWeight"`
	// Objectors are the peers voting against it, ordered by node id.
	Objectors     []ids.NodeID `json:"objectors"`
	AbstainWeight types.Uint64 `json:"abstainWeight"`
}

// LPStatus is one LP proposal and where the network's stake stands on it.
type LPStatus struct {
	// Number is the LP's number, the key it appears under on the JSON wire.
	Number uint32 `json:"number"`
	// LP is the stake for, against and abstaining.
	LP LP `json:"lp"`
}

// LPs is where the network stands on every LP: an object keyed by LP number on
// the JSON wire, a list in ascending number here.
type LPs []LPStatus

func (l LPs) MarshalJSON() ([]byte, error) {
	m := make(map[uint32]LP, len(l))
	for _, e := range l {
		m[e.Number] = e.LP
	}
	return json.Marshal(m)
}

func (l *LPs) UnmarshalJSON(b []byte) error {
	var m map[uint32]LP
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*l = make(LPs, 0, len(m))
	for number, lp := range m {
		*l = append(*l, LPStatus{Number: number, LP: lp})
	}
	slices.SortFunc(*l, func(x, y LPStatus) int { return cmpUint32(x.Number, y.Number) })
	return nil
}

// LPsReply are the results from calling LPs.
type LPsReply struct {
	LPs LPs `json:"lps"`
}

// Find is the status of one LP, and whether the reply carries it.
func (l LPs) Find(number uint32) (*LPStatus, bool) {
	for i := range l {
		if l[i].Number == number {
			return &l[i], true
		}
	}
	return nil, false
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

// VMAlias is one VM and the names it answers to.
type VMAlias struct {
	// VM is the VM's id.
	VM ids.ID `json:"vm"`
	// Aliases are the names it also answers to.
	Aliases []string `json:"aliases"`
}

// VMAliases are the VMs installed on a node: an object keyed by VM id on the
// JSON wire, a list ordered by that id here.
type VMAliases []VMAlias

func (v VMAliases) MarshalJSON() ([]byte, error) {
	m := make(map[ids.ID][]string, len(v))
	for _, e := range v {
		m[e.VM] = e.Aliases
	}
	return json.Marshal(m)
}

func (v *VMAliases) UnmarshalJSON(b []byte) error {
	var m map[ids.ID][]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*v = make(VMAliases, 0, len(m))
	for vm, aliases := range m {
		*v = append(*v, VMAlias{VM: vm, Aliases: aliases})
	}
	slices.SortFunc(*v, func(x, y VMAlias) int { return x.VM.Compare(y.VM) })
	return nil
}

// FxName is one feature extension and what it is called.
type FxName struct {
	// Fx is the extension's id.
	Fx ids.ID `json:"fx"`
	// Name is what it is called.
	Name string `json:"name"`
}

// FxNames are the feature extensions a node knows: an object keyed by id on the
// JSON wire, a list ordered by that id here.
type FxNames []FxName

func (f FxNames) MarshalJSON() ([]byte, error) {
	m := make(map[ids.ID]string, len(f))
	for _, e := range f {
		m[e.Fx] = e.Name
	}
	return json.Marshal(m)
}

func (f *FxNames) UnmarshalJSON(b []byte) error {
	var m map[ids.ID]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*f = make(FxNames, 0, len(m))
	for fx, name := range m {
		*f = append(*f, FxName{Fx: fx, Name: name})
	}
	slices.SortFunc(*f, func(x, y FxName) int { return x.Fx.Compare(y.Fx) })
	return nil
}

// GetVMsReply contains the response metadata for GetVMs.
type GetVMsReply struct {
	VMs VMAliases `json:"vms"`
	Fxs FxNames   `json:"fxs"`
}

// cmpString and cmpUint32 order the lists above. Written out rather than reached
// for from cmp so the ordering a reply is in is stated where the reply is.
func cmpString(x, y string) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	}
	return 0
}

func cmpUint32(x, y uint32) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	}
	return 0
}
