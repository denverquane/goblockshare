package blockchain

import (
	// "fmt"
	"testing"
)

//TODO Add more tests for empty strings, such as for usernames, transactions, channels, etc.!

func TestCreateNewChannel(t *testing.T) {
	_, err := CreateNewChannel(AuthTransaction{"", "", "", ""}, "ADMIN")

	if err == nil {
		t.Fail()
	}

	chain := MakeInitialChain([]UserPassPair{{"user", "pass"}}, "f")
	_, erre := SetChannelChain("", chain)

	if erre == nil {
		t.Fail()
	}

	//_, err := CreateNewChannel(AuthTransaction{"", "", "", ""}, "ADMIN")

}