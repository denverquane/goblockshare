package common

import (
	"encoding/json"
	"testing"
)

func TestSignableTransaction_GetHash(t *testing.T) {
	sign := SignableTransaction{}
	h1 := sign.GetHash()
	sign.TxID = "sdfgsdfgsd"
	h2 := sign.GetHash()

	if string(h1) != string(h2) {
		t.Fail()
	}
}

func TestJSONSignableTransaction_ConvertToSignable(t *testing.T) {
	addr := GenerateNewPersonalAddress()
	origin := addr.ConvertToOriginInfo()

	json := JSONSignableTransaction{origin, json.RawMessage(`{"precomputed": true}`), PUBLISH_TORRENT,
		origin.PubKeyX, origin.PubKeyY, "sdfg"}
	new := json.ConvertToSignable()

	if json.Origin.Address != new.Origin.Address || json.Origin.PubKeyX != new.Origin.PubKeyX || json.Origin.PubKeyY != new.Origin.PubKeyY {
		t.Fail()
	}

	if json.TxID != new.TxID || json.TransactionType != new.TransactionType || json.R != new.R || json.S != new.S {
		t.Fail()
	}
}
