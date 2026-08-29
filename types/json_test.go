// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package types

import (
	"encoding/json"
	"testing"

	"github.com/holiman/uint256"
)

// Every numeric at this boundary is carried as a QUOTED decimal, because
// JSON numbers are float64 and a uint64 above 2^53 does not survive one.
func TestQuotedDecimalRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   any
		want string
	}{
		{"uint16 max", Uint16(65535), `"65535"`},
		{"uint32 max", Uint32(4294967295), `"4294967295"`},
		{"uint64 max", Uint64(18446744073709551615), `"18446744073709551615"`},
		{"uint256 max", Uint256(*uint256.MustFromDecimal("115792089237316195423570985008687907853269984665640564039457584007913129639935")),
			`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`},
		{"uint256 one wei over uint64", Uint256(*uint256.MustFromDecimal("18446744073709551616")), `"18446744073709551616"`},
	} {
		got, err := json.Marshal(tc.in)
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.name, err)
		}
		if string(got) != tc.want {
			t.Errorf("%s: got %s want %s", tc.name, got, tc.want)
		}
	}
}

func TestUint256SurvivesWhatUint64Cannot(t *testing.T) {
	// A token balance with 18 decimals passes 2^64 at ~18.4 tokens.
	const wei = "1000000000000000000000" // 1000 tokens
	var u Uint256
	if err := json.Unmarshal([]byte(`"`+wei+`"`), &u); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v := uint256.Int(u); v.Dec() != wei {
		t.Errorf("got %s want %s", v.Dec(), wei)
	}
	var small Uint64
	if err := json.Unmarshal([]byte(`"`+wei+`"`), &small); err == nil {
		t.Error("Uint64 accepted a value it cannot hold; that is the bug Uint256 exists for")
	}
}

func TestNullLeavesTheValueAlone(t *testing.T) {
	u16, u32, u64, f := Uint16(7), Uint32(7), Uint64(7), Float64(7)
	u256 := Uint256(*uint256.NewInt(7))
	for _, p := range []json.Unmarshaler{&u16, &u32, &u64, &u256, &f} {
		if err := p.UnmarshalJSON([]byte(Null)); err != nil {
			t.Fatalf("null: %v", err)
		}
	}
	if u16 != 7 || u32 != 7 || u64 != 7 || f != 7 {
		t.Error("null overwrote a value")
	}
	if v := uint256.Int(u256); v.Dec() != "7" {
		t.Error("null overwrote the 256-bit value")
	}
}

func TestUnquotedIsAlsoAccepted(t *testing.T) {
	var u Uint64
	if err := json.Unmarshal([]byte(`12345`), &u); err != nil || u != 12345 {
		t.Errorf("bare number: got %d err %v", u, err)
	}
}
