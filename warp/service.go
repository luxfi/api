// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import "context"

// Service defines the warp IPC API contract.
type Service interface {
	PublishBlockchain(ctx context.Context, args *PublishBlockchainArgs) (*PublishBlockchainReply, error)
	UnpublishBlockchain(ctx context.Context, args *UnpublishBlockchainArgs) (*EmptyReply, error)
	GetPublishedBlockchains(ctx context.Context) (*GetPublishedBlockchainsReply, error)
}
