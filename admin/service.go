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
}
