// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package health

import "context"

// Service defines the health API contract.
type Service interface {
	Readiness(ctx context.Context, args *APIArgs) (*APIReply, error)
	Health(ctx context.Context, args *APIArgs) (*APIReply, error)
	Liveness(ctx context.Context, args *APIArgs) (*APIReply, error)
}
