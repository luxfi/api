// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package admin

import "context"

// Service defines the admin API contract.
type Service interface {
	StartCPUProfiler(ctx context.Context) (*EmptyReply, error)
	StopCPUProfiler(ctx context.Context) (*EmptyReply, error)
	MemoryProfile(ctx context.Context) (*EmptyReply, error)
	LockProfile(ctx context.Context) (*EmptyReply, error)

	Alias(ctx context.Context, args *AliasArgs) (*EmptyReply, error)
	AliasChain(ctx context.Context, args *AliasChainArgs) (*EmptyReply, error)
	GetChainAliases(ctx context.Context, args *GetChainAliasesArgs) (*GetChainAliasesReply, error)

	Stacktrace(ctx context.Context) (*EmptyReply, error)
	SetLoggerLevel(ctx context.Context, args *SetLoggerLevelArgs) (*EmptyReply, error)
	GetLoggerLevel(ctx context.Context, args *GetLoggerLevelArgs) (*LoggerLevelReply, error)
	GetConfig(ctx context.Context) (any, error)

	LoadVMs(ctx context.Context) (*LoadVMsReply, error)
	DbGet(ctx context.Context, args *DBGetArgs) (*DBGetReply, error)
	ListVMs(ctx context.Context) (*ListVMsReply, error)

	SetTrackedChains(ctx context.Context, args *SetTrackedChainsArgs) (*SetTrackedChainsReply, error)
	GetTrackedChains(ctx context.Context) (*GetTrackedChainsReply, error)

	// Snapshot creates a database backup to the specified path.
	// Use .zst extension for zstd compression.
	// since: 0 for full backup, or version from last backup for incremental.
	// Returns the version number for future incremental backups.
	Snapshot(ctx context.Context, args *SnapshotArgs) (*SnapshotReply, error)

	// Load restores the database from a backup file.
	// Detects compression from .zst extension.
	Load(ctx context.Context, args *LoadArgs) (*EmptyReply, error)
}
