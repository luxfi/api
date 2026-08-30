// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// The shapes that changed to cross a wire must not have changed the wire.
//
// A map became a list and an `any` became bytes so that a reply could cross
// between two processes at all — a ZAP field is an offset, and a map has no
// order to give one while an interface has no shape until it holds something.
// None of that is a caller's business: these APIs are answered by a live network
// today, so the JSON a caller reads must be the same bytes it always was.
//
// Each test below writes the NEW type and compares against the object the OLD
// one produced, spelled out as a literal. A literal rather than a round trip,
// because a round trip through one type proves only that it agrees with itself.
package api_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/luxfi/api/admin"
	"github.com/luxfi/api/health"
	"github.com/luxfi/api/info"
	"github.com/luxfi/api/types"
	"github.com/luxfi/ids"
)

func id(b byte) ids.ID { return ids.ID{0: b} }

// same asserts v marshals to want, and that reading want back yields v — so the
// list is not merely writing the old shape but reading it too.
func same[T any](t *testing.T, v T, want string, back func(T) bool) {
	t.Helper()
	got, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if norm(t, string(got)) != norm(t, want) {
		t.Errorf("marshalled\n  %s\nwant\n  %s", got, want)
	}
	var read T
	if err := json.Unmarshal([]byte(want), &read); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !back(read) {
		t.Errorf("reading %s back did not reproduce the value", want)
	}
}

func norm(t *testing.T, s string) string {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatalf("not JSON: %s", s)
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func TestVMVersionsIsStillAnObject(t *testing.T) {
	same(t,
		info.VMVersions{{VM: "platformvm", Version: "v1.2.3"}, {VM: "xvm", Version: "v1.2.4"}},
		`{"platformvm":"v1.2.3","xvm":"v1.2.4"}`,
		func(v info.VMVersions) bool {
			return len(v) == 2 && v[0].VM == "platformvm" && v[1].Version == "v1.2.4"
		})
}

func TestVMAliasesAndFxNamesAreStillObjects(t *testing.T) {
	same(t,
		info.VMAliases{{VM: id(1), Aliases: []string{"platformvm"}}},
		`{"`+id(1).String()+`":["platformvm"]}`,
		func(v info.VMAliases) bool { return len(v) == 1 && v[0].VM == id(1) })

	same(t,
		info.FxNames{{Fx: id(2), Name: "secp256k1fx"}},
		`{"`+id(2).String()+`":"secp256k1fx"}`,
		func(v info.FxNames) bool { return len(v) == 1 && v[0].Name == "secp256k1fx" })
}

func TestLPsAreStillAnObjectKeyedByNumber(t *testing.T) {
	same(t,
		info.LPs{{Number: 23, LP: info.LP{AbstainWeight: 7}}},
		`{"23":{"supportWeight":"0","supporters":null,"objectWeight":"0","objectors":null,"abstainWeight":"7"}}`,
		func(v info.LPs) bool { return len(v) == 1 && v[0].Number == 23 && v[0].LP.AbstainWeight == 7 })
}

func TestLoggerLevelsAreStillAnObject(t *testing.T) {
	same(t,
		admin.LoggerLevels{{Logger: "C", Levels: admin.LogAndDisplayLevels{LogLevel: "DEBUG", DisplayLevel: "INFO"}}},
		`{"C":{"logLevel":"DEBUG","displayLevel":"INFO"}}`,
		func(v admin.LoggerLevels) bool { return len(v) == 1 && v[0].Levels.LogLevel == "DEBUG" })
}

func TestInstalledVMsAreStillAnObjectKeyedByID(t *testing.T) {
	same(t,
		admin.InstalledVMs{{ID: "vm1", Aliases: []string{"a"}}},
		`{"vm1":{"id":"vm1","aliases":["a"]}}`,
		func(v admin.InstalledVMs) bool { return len(v) == 1 && v[0].ID == "vm1" })
}

func TestLoadedAndFailedVMsAreStillObjects(t *testing.T) {
	same(t,
		admin.LoadedVMs{{VM: id(3), Aliases: []string{"a"}}},
		`{"`+id(3).String()+`":["a"]}`,
		func(v admin.LoadedVMs) bool { return len(v) == 1 && v[0].VM == id(3) })

	same(t,
		admin.FailedVMs{{VM: id(4), Error: "no plugin"}},
		`{"`+id(4).String()+`":"no plugin"}`,
		func(v admin.FailedVMs) bool { return len(v) == 1 && v[0].Error == "no plugin" })
}

func TestHealthChecksAreStillAnObject(t *testing.T) {
	when := types.TimeOf(time.Unix(1753479996, 0).UTC())
	said, err := types.RawOf(map[string]any{"availableDiskBytes": 12})
	if err != nil {
		t.Fatal(err)
	}
	same(t,
		health.Checks{{Name: "diskspace", Result: health.Result{Details: said, Timestamp: when, Duration: 1234}}},
		`{"diskspace":{"message":{"availableDiskBytes":12},"timestamp":"2025-07-25T21:46:36Z","duration":1234}}`,
		func(v health.Checks) bool { return len(v) == 1 && v[0].Name == "diskspace" })
}

// The details of a check are whatever it said. All three shapes mainnet answers
// with survive as themselves — an `any` that became bytes lost nothing.
func TestRawKeepsEveryShapeACheckAnswersWith(t *testing.T) {
	for _, want := range []string{
		`{"availableDiskBytes":12}`,
		`"node is not a validator"`,
		`["P","X","C"]`,
		`null`,
	} {
		var r types.Raw
		if err := json.Unmarshal([]byte(want), &r); err != nil {
			t.Fatalf("unmarshal %s: %v", want, err)
		}
		got, err := json.Marshal(r)
		if err != nil {
			t.Fatalf("marshal %s: %v", want, err)
		}
		if string(got) != want {
			t.Errorf("Raw turned %s into %s", want, got)
		}
	}
}

// An address is text on the JSON wire and always was — including the empty
// string an invalid one has always been written as.
func TestAddrIsTheTextNetipWrites(t *testing.T) {
	for _, want := range []string{`"203.0.113.9:9651"`, `""`} {
		var a types.Addr
		if err := json.Unmarshal([]byte(want), &a); err != nil {
			t.Fatalf("unmarshal %s: %v", want, err)
		}
		got, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("Addr turned %s into %s", want, got)
		}
	}
}

// An instant is the RFC 3339 string time.Time writes, and reads back to the
// same instant in the same zone.
func TestTimeIsTheStringTimeTimeWrites(t *testing.T) {
	when := time.Date(2026, 8, 29, 12, 0, 0, 0, time.FixedZone("", 2*3600))
	got, err := json.Marshal(types.TimeOf(when))
	if err != nil {
		t.Fatal(err)
	}
	direct, err := json.Marshal(when)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(direct) {
		t.Errorf("Time wrote %s; time.Time writes %s", got, direct)
	}

	var read types.Time
	if err := json.Unmarshal(direct, &read); err != nil {
		t.Fatal(err)
	}
	if !read.Time().Equal(when) {
		t.Errorf("read back %s; want %s", read.Time(), when)
	}
}
