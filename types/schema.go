// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package types

// Each numeric here writes itself as a quoted decimal, so a schema built from
// its fields would describe the Go integer and not the wire. Stating the shape
// beside MarshalJSON is what keeps the two from disagreeing, and it is what
// lets a generated client type the field instead of receiving an untyped
// value. zip reads these through its SchemaDescriber interface.
//
// The pattern is the constraint that matters: the value is digits in a string,
// which is the whole reason these types exist — a JSON number is a float64,
// and a uint64 above 2^53 does not survive one.

func quotedInteger(bits string) map[string]any {
	return map[string]any{
		"type":        "string",
		"pattern":     "^[0-9]+$",
		"description": "An unsigned " + bits + "-bit integer, carried as a decimal string.",
	}
}

func (Uint16) JSONSchema() map[string]any { return quotedInteger("16") }
func (Uint32) JSONSchema() map[string]any { return quotedInteger("32") }
func (Uint64) JSONSchema() map[string]any { return quotedInteger("64") }

// Uint256 is the width the EVM carries; it does not fit any JSON numeric.
func (Uint256) JSONSchema() map[string]any { return quotedInteger("256") }

func (Float64) JSONSchema() map[string]any {
	return map[string]any{
		"type":        "string",
		"pattern":     "^-?[0-9]+\\.[0-9]+$",
		"description": "A 64-bit float, carried as a decimal string.",
	}
}
