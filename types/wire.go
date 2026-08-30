// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Values that are one thing at the JSON edge and another on the ZAP wire.
//
// A ZAP field is an OFFSET and a WIDTH. Three shapes a Go program reaches for
// have neither, and each fails differently:
//
//   - a MAP has no order, so it has no layout at all. [zip.LayoutOf] refuses it
//     outright, which is the loud failure and the easy one.
//   - a struct with only UNEXPORTED fields — time.Time, netip.AddrPort —
//     derives an EMPTY layout instead. It is not refused: the value crosses,
//     reports no error and carries nothing, which is the quiet failure and the
//     dangerous one.
//   - an INTERFACE has no shape to lay out until it holds something, and what it
//     holds is a runtime fact.
//
// The types here are the second and third cases answered. The first is answered
// per map, in the package that owns it: an object on the JSON wire, a list
// ordered by key here.
//
// EVERY ONE OF THEM RENDERS THE JSON IT REPLACED, byte for byte. That is the
// whole constraint — these travel on APIs a network already answers, so the
// change is which bytes cross between two processes, never which bytes reach a
// caller.

package types

import (
	"encoding/json"
	"net/netip"
	"time"
)

// Addr is a network address — an IP and a port — as the JSON wire has always
// carried it: the text netip.AddrPort writes for itself.
//
// It is text here rather than a netip.AddrPort because that type keeps its
// address, its zone and its port in unexported fields. Reflection cannot read
// them, so the derived layout is empty and a reply holding one arrives blank
// with no error to say so. GetNodeIPReply has exactly one field, so its entire
// answer was empty.
type Addr string

// AddrOf is ap as it crosses. It renders through netip's own marshaler, so the
// bytes are the ones that were already on the JSON wire — including the empty
// string an invalid address has always been written as.
func AddrOf(ap netip.AddrPort) Addr {
	b, err := ap.MarshalText()
	if err != nil {
		return ""
	}
	return Addr(b)
}

// AddrPort is the address as a netip.AddrPort, or the zero value when the text
// does not parse — which is what an absent address already reads as.
func (a Addr) AddrPort() netip.AddrPort {
	ap, err := netip.ParseAddrPort(string(a))
	if err != nil {
		return netip.AddrPort{}
	}
	return ap
}

// Time is an instant, as three numbers that have a layout.
//
// time.Time cannot cross for the same reason netip.AddrPort cannot: every field
// is unexported. The JSON is byte for byte what time.Time's own marshaler
// writes, because the rendering depends only on the instant and the offset —
// RFC 3339 spells a zero offset "Z" and any other one "+HH:MM" — so carrying
// the offset is enough to put the instant back in the zone it was read in.
type Time struct {
	// Seconds is seconds since the epoch.
	Seconds int64 `json:"seconds"`
	// Nanos is the nanosecond within that second.
	Nanos int32 `json:"nanos"`
	// Offset is seconds east of UTC.
	Offset int32 `json:"offset"`
}

// TimeOf is t as a value that can cross.
func TimeOf(t time.Time) Time {
	_, offset := t.Zone()
	return Time{Seconds: t.Unix(), Nanos: int32(t.Nanosecond()), Offset: int32(offset)}
}

// Time is the instant, in the zone it was read in.
func (t Time) Time() time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos)).In(time.FixedZone("", int(t.Offset)))
}

// IsZero reports whether this is the zero instant, which is what an absent time
// has always marshalled as.
func (t Time) IsZero() bool { return t == Time{} }

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal(time.Time{})
	}
	return json.Marshal(t.Time())
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if string(b) == Null {
		return nil
	}
	var when time.Time
	if err := json.Unmarshal(b, &when); err != nil {
		return err
	}
	*t = TimeOf(when)
	return nil
}

// Raw is a value whose shape only the sender knows: it crosses as the bytes of
// its own JSON, and is written back out as those bytes.
//
// It is what an `any` field becomes. An interface has no layout until it holds
// something, and what it holds is decided per call — a health check answers with
// an object, a string or an array depending on which check it is — so there is
// no struct to name. The bytes are the honest shape, and they reach a caller
// exactly as they were produced.
type Raw []byte

// RawOf is v encoded once, at the point where its concrete type is still known.
func RawOf(v any) (Raw, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return Raw(b), nil
}

func (r Raw) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return []byte(Null), nil
	}
	return r, nil
}

func (r *Raw) UnmarshalJSON(b []byte) error {
	if string(b) == Null {
		*r = nil
		return nil
	}
	*r = append((*r)[:0], b...)
	return nil
}
