// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Operating a node.
//
// THREE SHAPES HERE ARE AN OBJECT ON THE JSON WIRE AND A LIST IN GO — the logger
// levels, the installed VMs, and what a reload found. A ZAP field is an offset
// and a map has no order to give one, so a reply carrying a map cannot cross
// between two processes at all. Each list is ordered by the key it was an object
// of; the JSON is unchanged.
package admin

import (
	"encoding/json"
	"slices"
	"strings"

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

// LoggerLevel is one logger and the two levels it is set to.
type LoggerLevel struct {
	// Logger is the logger's name, the key it appears under on the JSON wire.
	Logger string `json:"logger"`
	// Levels are what it writes to its file and to the display.
	Levels LogAndDisplayLevels `json:"levels"`
}

// LoggerLevels are the levels every named logger is set to: an object keyed by
// logger name on the JSON wire, a list ordered by that name here.
type LoggerLevels []LoggerLevel

func (l LoggerLevels) MarshalJSON() ([]byte, error) {
	m := make(map[string]LogAndDisplayLevels, len(l))
	for _, e := range l {
		m[e.Logger] = e.Levels
	}
	return json.Marshal(m)
}

func (l *LoggerLevels) UnmarshalJSON(b []byte) error {
	var m map[string]LogAndDisplayLevels
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*l = make(LoggerLevels, 0, len(m))
	for name, levels := range m {
		*l = append(*l, LoggerLevel{Logger: name, Levels: levels})
	}
	slices.SortFunc(*l, func(x, y LoggerLevel) int { return strings.Compare(x.Logger, y.Logger) })
	return nil
}

// LoggerLevelReply are the levels of the loggers.
type LoggerLevelReply struct {
	LoggerLevels LoggerLevels `json:"loggerLevels"`
}

// GetLoggerLevelArgs are the arguments for getting logger levels.
type GetLoggerLevelArgs struct {
	LoggerName string `json:"loggerName"`
}

// LoadedVM is one VM a reload brought in and the names it answers to.
type LoadedVM struct {
	// VM is the VM's id, the key it appears under on the JSON wire.
	VM ids.ID `json:"vm"`
	// Aliases are the names it also answers to.
	Aliases []string `json:"aliases"`
}

// LoadedVMs are the VMs a reload brought in: an object keyed by VM id on the
// JSON wire, a list ordered by that id here.
type LoadedVMs []LoadedVM

func (l LoadedVMs) MarshalJSON() ([]byte, error) {
	m := make(map[ids.ID][]string, len(l))
	for _, e := range l {
		m[e.VM] = e.Aliases
	}
	return json.Marshal(m)
}

func (l *LoadedVMs) UnmarshalJSON(b []byte) error {
	var m map[ids.ID][]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*l = make(LoadedVMs, 0, len(m))
	for vm, aliases := range m {
		*l = append(*l, LoadedVM{VM: vm, Aliases: aliases})
	}
	slices.SortFunc(*l, func(x, y LoadedVM) int { return x.VM.Compare(y.VM) })
	return nil
}

// FailedVM is one VM a reload could not bring in, and why.
type FailedVM struct {
	// VM is the VM's id, the key it appears under on the JSON wire.
	VM ids.ID `json:"vm"`
	// Error is what stopped it loading.
	Error string `json:"error"`
}

// FailedVMs are the VMs a reload could not bring in: an object keyed by VM id
// on the JSON wire, a list ordered by that id here.
type FailedVMs []FailedVM

func (f FailedVMs) MarshalJSON() ([]byte, error) {
	m := make(map[ids.ID]string, len(f))
	for _, e := range f {
		m[e.VM] = e.Error
	}
	return json.Marshal(m)
}

func (f *FailedVMs) UnmarshalJSON(b []byte) error {
	var m map[ids.ID]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*f = make(FailedVMs, 0, len(m))
	for vm, why := range m {
		*f = append(*f, FailedVM{VM: vm, Error: why})
	}
	slices.SortFunc(*f, func(x, y FailedVM) int { return x.VM.Compare(y.VM) })
	return nil
}

// LoadVMsReply contains the response for LoadVMs.
type LoadVMsReply struct {
	NewVMs        LoadedVMs `json:"newVMs"`
	FailedVMs     FailedVMs `json:"failedVMs"`
	ChainsRetried int       `json:"chainsRetried"`
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

// InstalledVMs are the VMs on a node: an object keyed by VM id on the JSON
// wire, a list ordered by that id here. The id is on the element too — it is
// [VMInfo.ID] — so the key is the element's own name rather than a second one.
type InstalledVMs []VMInfo

func (v InstalledVMs) MarshalJSON() ([]byte, error) {
	m := make(map[string]VMInfo, len(v))
	for _, e := range v {
		m[e.ID] = e
	}
	return json.Marshal(m)
}

func (v *InstalledVMs) UnmarshalJSON(b []byte) error {
	var m map[string]VMInfo
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*v = make(InstalledVMs, 0, len(m))
	for id, info := range m {
		info.ID = id
		*v = append(*v, info)
	}
	slices.SortFunc(*v, func(x, y VMInfo) int { return strings.Compare(x.ID, y.ID) })
	return nil
}

// ListVMsReply contains the response for ListVMs.
type ListVMsReply struct {
	VMs InstalledVMs `json:"vms"`
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
