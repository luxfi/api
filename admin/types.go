// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package admin

import (
	"github.com/luxfi/api/types"
	"github.com/luxfi/ids"
)

// AliasArgs are the arguments for calling Alias.
type AliasArgs struct {
	Endpoint string `json:"endpoint"`
	Alias    string `json:"alias"`
}

// AliasChainArgs are the arguments for calling AliasChain.
type AliasChainArgs struct {
	Chain string `json:"chain"`
	Alias string `json:"alias"`
}

// GetChainAliasesArgs are the arguments for calling GetChainAliases.
type GetChainAliasesArgs struct {
	Chain string `json:"chain"`
}

// GetChainAliasesReply are the aliases of the given chain.
type GetChainAliasesReply struct {
	Aliases []string `json:"aliases"`
}

// SetLoggerLevelArgs are the arguments for setting a logger's levels.
type SetLoggerLevelArgs struct {
	LoggerName   string `json:"loggerName"`
	LogLevel     string `json:"logLevel"`
	DisplayLevel string `json:"displayLevel"`
}

// LogAndDisplayLevels pairs log and display levels.
type LogAndDisplayLevels struct {
	LogLevel     string `json:"logLevel"`
	DisplayLevel string `json:"displayLevel"`
}

// LoggerLevelReply are the levels of the loggers.
type LoggerLevelReply struct {
	LoggerLevels map[string]LogAndDisplayLevels `json:"loggerLevels"`
}

// GetLoggerLevelArgs are the arguments for getting logger levels.
type GetLoggerLevelArgs struct {
	LoggerName string `json:"loggerName"`
}

// LoadVMsReply contains the response for LoadVMs.
type LoadVMsReply struct {
	NewVMs        map[ids.ID][]string `json:"newVMs"`
	FailedVMs     map[ids.ID]string   `json:"failedVMs"`
	ChainsRetried int                 `json:"chainsRetried"`
}

// DBGetArgs are the arguments for DBGet.
type DBGetArgs struct {
	Key string `json:"key"`
}

// DBGetReply is the reply for DBGet.
type DBGetReply struct {
	Value string `json:"value"`
}

// VMInfo contains information about a registered VM.
type VMInfo struct {
	ID      string   `json:"id"`
	Aliases []string `json:"aliases"`
	Path    string   `json:"path,omitempty"`
}

// ListVMsReply contains the response for ListVMs.
type ListVMsReply struct {
	VMs map[string]VMInfo `json:"vms"`
}

// SnapshotArgs are the arguments for Snapshot.
type SnapshotArgs struct {
	Path  string `json:"path"`
	Since uint64 `json:"since"`
}

// SnapshotReply is the response for Snapshot.
type SnapshotReply struct {
	Version uint64 `json:"version"`
}

// LoadArgs are the arguments for Load.
type LoadArgs struct {
	Path string `json:"path"`
}

// SetTrackedChainsArgs are the arguments for SetTrackedChains.
type SetTrackedChainsArgs struct {
	Chains []string `json:"chains"`
}

// SetTrackedChainsReply is the response for SetTrackedChains.
type SetTrackedChainsReply struct {
	TrackedChains []string `json:"trackedChains"`
}

// GetTrackedChainsReply is the response for GetTrackedChains.
type GetTrackedChainsReply struct {
	TrackedChains []string `json:"trackedChains"`
}

// EmptyReply is an alias to common empty reply.
type EmptyReply = types.EmptyReply
