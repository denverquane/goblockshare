package blockchain

import "testing"

//TODO Add more tests!

func TestCreateNewChannel(t *testing.T) {
	trans := AuthTransaction{"", "", "", ""}
	_, err := CreateNewGlobalChannel(trans)

	if err == nil {
		t.Fail()
	}

	chain := BlockChain{}

	_, err2 := SetChannelChain("", chain) //valid

	if err2 == nil {
		t.Fail()
	}

	_, err3 := SetChannelChain("ADMIN", chain) //valid name, but empty chain
	if err3 == nil {
		t.Fail()
	}
}
