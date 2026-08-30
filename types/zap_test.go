// Copyright (C) 2019-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luxfi/formatting"
	"github.com/luxfi/ids"
)

// These are the shapes every chain's RPC answers in, and each of them holds an
// ids.ID — [32]byte, which is bytes_fixed[32] on the ZAP wire and what a DERIVED
// codec refuses to carry. zap_gen.go is why they cross; this is that the id
// arrives whole and that the JSON mainnet answers on did not move.
func TestTheIDsCrossAndTheJSONDoesNotMove(t *testing.T) {
	id := ids.ID{0: 0xf0, 15: 0x5a, 31: 0x0d}

	t.Run("JSONTxID", func(t *testing.T) {
		require := require.New(t)
		sent := JSONTxID{TxID: id}
		enc, err := sent.MarshalZAP()
		require.NoError(err)
		var back JSONTxID
		require.NoError(back.UnmarshalZAP(enc))
		require.Equal(sent, back)

		is, err := json.Marshal(sent)
		require.NoError(err)
		was, err := json.Marshal(struct {
			TxID ids.ID `json:"txID"`
		}{id})
		require.NoError(err)
		require.Equal(string(was), string(is))
	})

	t.Run("GetTxArgs", func(t *testing.T) {
		require := require.New(t)
		sent := GetTxArgs{TxID: id, Encoding: formatting.Hex}
		enc, err := sent.MarshalZAP()
		require.NoError(err)
		var back GetTxArgs
		require.NoError(back.UnmarshalZAP(enc))
		require.Equal(sent, back)
	})

	t.Run("GetBlockArgs", func(t *testing.T) {
		require := require.New(t)
		sent := GetBlockArgs{BlockID: id, Encoding: formatting.JSON}
		enc, err := sent.MarshalZAP()
		require.NoError(err)
		var back GetBlockArgs
		require.NoError(back.UnmarshalZAP(enc))
		require.Equal(sent, back)
	})

	t.Run("GetUTXOsArgs", func(t *testing.T) {
		require := require.New(t)
		sent := GetUTXOsArgs{
			Addresses:   []string{"P-a", "P-b"},
			SourceChain: "X",
			Limit:       10,
			Encoding:    formatting.Hex,
		}
		sent.StartIndex.Address = "P-a"
		sent.StartIndex.UTXO = "utxo"
		enc, err := sent.MarshalZAP()
		require.NoError(err)
		var back GetUTXOsArgs
		require.NoError(back.UnmarshalZAP(enc))
		require.Equal(sent, back)
	})
}
