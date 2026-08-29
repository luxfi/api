// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package types

import (
	"strconv"

	"github.com/holiman/uint256"
)

const Null = "null"

// payload strips the quotes JSON-RPC numerics are carried in and says whether
// the value was null. Every type below reads a quoted decimal, so the reading
// is stated once; each type differs only in how it parses what comes out.
func payload(b []byte) (string, bool) {
	s := string(b)
	if s == Null {
		return "", true
	}
	if len(s) >= 2 {
		if last := len(s) - 1; s[0] == '"' && s[last] == '"' {
			s = s[1:last]
		}
	}
	return s, false
}

func quote(s string) []byte { return []byte(`"` + s + `"`) }

type Uint16 uint16

func (u Uint16) MarshalJSON() ([]byte, error) {
	return quote(strconv.FormatUint(uint64(u), 10)), nil
}

func (u *Uint16) UnmarshalJSON(b []byte) error {
	s, null := payload(b)
	if null {
		return nil
	}
	val, err := strconv.ParseUint(s, 10, 16)
	*u = Uint16(val)
	return err
}

type Uint32 uint32

func (u Uint32) MarshalJSON() ([]byte, error) {
	return quote(strconv.FormatUint(uint64(u), 10)), nil
}

func (u *Uint32) UnmarshalJSON(b []byte) error {
	s, null := payload(b)
	if null {
		return nil
	}
	val, err := strconv.ParseUint(s, 10, 32)
	*u = Uint32(val)
	return err
}

type Uint64 uint64

func (u Uint64) MarshalJSON() ([]byte, error) {
	return quote(strconv.FormatUint(uint64(u), 10)), nil
}

func (u *Uint64) UnmarshalJSON(b []byte) error {
	s, null := payload(b)
	if null {
		return nil
	}
	val, err := strconv.ParseUint(s, 10, 64)
	*u = Uint64(val)
	return err
}

// Uint256 is the width the EVM actually carries — a balance, a wei value or a
// token amount does not fit in 64 bits, and every such field reaching this
// boundary had to be a string or a big.Int by hand. It is fixed-size rather
// than big.Int so it cannot be nil and cannot allocate on the hot path.
type Uint256 uint256.Int

func (u Uint256) MarshalJSON() ([]byte, error) {
	v := uint256.Int(u)
	return quote(v.Dec()), nil
}

func (u *Uint256) UnmarshalJSON(b []byte) error {
	s, null := payload(b)
	if null {
		return nil
	}
	v, err := uint256.FromDecimal(s)
	if err != nil {
		return err
	}
	*u = Uint256(*v)
	return nil
}

type Float64 float64

func (f Float64) MarshalJSON() ([]byte, error) {
	return quote(strconv.FormatFloat(float64(f), 'f', 4, 64)), nil
}

func (f *Float64) UnmarshalJSON(b []byte) error {
	s, null := payload(b)
	if null {
		return nil
	}
	val, err := strconv.ParseFloat(s, 64)
	*f = Float64(val)
	return err
}
