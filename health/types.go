// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Whether a node is up, ready, and well.
//
// TWO SHAPES HERE COULD NOT CROSS BETWEEN TWO PROCESSES and now do. The checks
// were a map, and a ZAP field is an offset a map has no order to give; they are
// a list ordered by check name. And a check's details were an `any`, which has
// no shape until it holds something — the disk check answers with an object, the
// BLS check with a string, and a chain's check with whatever it likes — so the
// details cross as the bytes of their own JSON, which is the one shape they all
// share. Neither change is visible to a caller: the list marshals to the object
// it replaced, and the raw bytes are written out exactly as they were produced.
package health

import (
	"encoding/json"
	"slices"
	"strings"
	"time"

	"github.com/luxfi/api/types"
)

// APIReply is the response for Readiness, Health, and Liveness.
type APIReply struct {
	// Checks is what every check answered, ordered by name.
	Checks Checks `json:"checks"`
	// Healthy is whether every one of them passed.
	Healthy bool `json:"healthy"`
}

// StatusCode is the HTTP status this answer carries: 200 when every check
// passed, 503 when one did not.
//
// It rides the value rather than being set beside it so the code a probe reads
// and the `healthy` a person reads are the same fact computed once. A caller
// that never speaks HTTP is unaffected — the value is identical on every wire.
func (r APIReply) StatusCode() int {
	if r.Healthy {
		return 200
	}
	return 503
}

// APIArgs is the arguments for Readiness, Health, and Liveness.
type APIArgs struct {
	// Tags narrows the answer to checks carrying one of these tags. Empty asks
	// for all of them.
	Tags []string `json:"tags"`
}

// Check is one named check and what it answered.
type Check struct {
	// Name is the check's name, the key it appears under on the JSON wire.
	Name string `json:"name"`
	// Result is what it answered.
	Result Result `json:"result"`
}

// Checks is what every check answered: an object keyed by check name on the
// JSON wire, a list ordered by that name here.
type Checks []Check

func (c Checks) MarshalJSON() ([]byte, error) {
	m := make(map[string]Result, len(c))
	for _, e := range c {
		m[e.Name] = e.Result
	}
	return json.Marshal(m)
}

func (c *Checks) UnmarshalJSON(b []byte) error {
	var m map[string]Result
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*c = make(Checks, 0, len(m))
	for name, result := range m {
		*c = append(*c, Check{Name: name, Result: result})
	}
	slices.SortFunc(*c, func(x, y Check) int { return strings.Compare(x.Name, y.Name) })
	return nil
}

// Find is one check's result, and whether the reply carries it.
func (c Checks) Find(name string) (*Check, bool) {
	for i := range c {
		if c[i].Name == name {
			return &c[i], true
		}
	}
	return nil, false
}

// Result describes a health check result.
type Result struct {
	// Details is whatever the check chose to say, as the JSON it said it in.
	// A check answers with an object, a string or a list as it sees fit, so
	// there is no struct to name and the bytes are the honest shape.
	Details types.Raw `json:"message,omitempty"`
	// Error is why the check failed, when it did.
	Error *string `json:"error,omitempty"`
	// Timestamp is when the check last ran.
	Timestamp types.Time `json:"timestamp,omitempty"`
	// Duration is how long it took.
	Duration time.Duration `json:"duration"`
	// ContiguousFailures is how many times running it has failed.
	ContiguousFailures int64 `json:"contiguousFailures,omitempty"`
	// TimeOfFirstFailure is when the current run of failures began.
	TimeOfFirstFailure *types.Time `json:"timeOfFirstFailure,omitempty"`
}
