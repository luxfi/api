// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package zap provides Zero-Copy App Proto (ZAP) serialization for high-performance
// RPC communication. ZAP minimizes CPU overhead and memory allocations through:
//
//   - Direct memory access for fixed-size fields
//   - Length-prefixed variable data without copying
//   - Buffer pooling to avoid allocations
//   - Simple framing for TCP transport (no HTTP/2 overhead)
//
// Wire protocol:
//
//	[4 bytes: message length][1 byte: message type][payload...]
//
// The message type allows multiplexing multiple RPC calls over a single connection.
package zap
