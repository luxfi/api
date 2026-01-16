// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package health

import "time"

// APIReply is the response for Readiness, Health, and Liveness.
type APIReply struct {
	Checks  map[string]Result `json:"checks"`
	Healthy bool              `json:"healthy"`
}

// APIArgs is the arguments for Readiness, Health, and Liveness.
type APIArgs struct {
	Tags []string `json:"tags"`
}

// Result describes a health check result.
type Result struct {
	Details             interface{} `json:"message,omitempty"`
	Error               *string     `json:"error,omitempty"`
	Timestamp           time.Time   `json:"timestamp,omitempty"`
	Duration            time.Duration `json:"duration"`
	ContiguousFailures  int64       `json:"contiguousFailures,omitempty"`
	TimeOfFirstFailure  *time.Time  `json:"timeOfFirstFailure,omitempty"`
}
