package common

import "testing"

func TestGenerateNewPersonalAddress(t *testing.T) {
	addr1 := GenerateNewPersonalAddress()
	addr2 := GenerateNewPersonalAddress()

	if addr1.Address == addr2.Address {
		t.Fail()
	}
}

func TestPersonalAddress_ConvertToOriginInfo(t *testing.T) {
	addr := GenerateNewPersonalAddress()
	origin := addr.ConvertToOriginInfo()

	if addr.Address != origin.Address || addr.PublicKey.X != origin.PubKeyX || addr.PublicKey.Y != origin.PubKeyY {
		t.Fail()
	}
}