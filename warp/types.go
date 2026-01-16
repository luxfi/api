// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"github.com/luxfi/api/types"
	"github.com/luxfi/ids"
)

type PublishBlockchainArgs struct {
	BlockchainID string `json:"blockchainID"`
}

type PublishBlockchainReply struct {
	ConsensusURL string `json:"consensusURL"`
	DecisionsURL string `json:"decisionsURL"`
}

type UnpublishBlockchainArgs struct {
	BlockchainID string `json:"blockchainID"`
}

type GetPublishedBlockchainsReply struct {
	Chains []ids.ID `json:"chains"`
}

type EmptyReply = types.EmptyReply
