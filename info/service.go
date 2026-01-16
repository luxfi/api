// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package info

import "context"

const (
	MethodGetNodeVersion = "info.getNodeVersion"
	MethodGetNodeID      = "info.getNodeID"
	MethodGetNodeIP      = "info.getNodeIP"
	MethodGetNetworkID   = "info.getNetworkID"
	MethodGetNetworkName = "info.getNetworkName"
	MethodGetBlockchainID = "info.getBlockchainID"
	MethodPeers          = "info.peers"
	MethodIsBootstrapped = "info.isBootstrapped"
	MethodUpgrades       = "info.upgrades"
	MethodUptime         = "info.uptime"
	MethodLps            = "info.lps"
	MethodGetTxFee       = "info.getTxFee"
	MethodGetVMs         = "info.getVMs"
)

// Service defines the info API contract.
type Service interface {
	GetNodeVersion(ctx context.Context) (*GetNodeVersionReply, error)
	GetNodeID(ctx context.Context) (*GetNodeIDReply, error)
	GetNodeIP(ctx context.Context) (*GetNodeIPReply, error)
	GetNetworkID(ctx context.Context) (*GetNetworkIDReply, error)
	GetNetworkName(ctx context.Context) (*GetNetworkNameReply, error)
	GetBlockchainID(ctx context.Context, args *GetBlockchainIDArgs) (*GetBlockchainIDReply, error)
	Peers(ctx context.Context, args *PeersArgs) (*PeersReply, error)
	IsBootstrapped(ctx context.Context, args *IsBootstrappedArgs) (*IsBootstrappedResponse, error)
	Upgrades(ctx context.Context) (*map[string]any, error)
	Uptime(ctx context.Context) (*UptimeResponse, error)
	Lps(ctx context.Context) (*LPsReply, error)
	GetTxFee(ctx context.Context) (*GetTxFeeResponse, error)
	GetVMs(ctx context.Context) (*GetVMsReply, error)
}
