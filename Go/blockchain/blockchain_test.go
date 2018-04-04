package blockchain

import (
	"fmt"
	"testing"
)

func TestMakeInitialChain(t *testing.T) {
	chain := MakeInitialChain([]UserPassPair{}, "")

	fmt.Println(chain.ToString())

	if len(chain.Blocks) != 1 {
		t.Fail()
	}

	if len(chain.Blocks[0].Users) != 0 {
		t.Fail()
	}

	if len(chain.Blocks[0].Transactions) != 0 {
		t.Fail()
	}

	if chain.Blocks[0].Index != 0 {
		t.Fail()
	}

}

func TestCreateChainFromSeed(t *testing.T) {
	chain := MakeInitialChain([]UserPassPair{}, "v")
	newChain := CreateChainFromSeed(chain)

	if len(newChain.Blocks) != 1 {
		t.Fail()
	}

	if len(newChain.Blocks[0].Users) != 0 {
		t.Fail()
	}

	chain2 := MakeInitialChain([]UserPassPair{{"user", "pass"}}, "")
	newChain2 := CreateChainFromSeed(chain2)

	if len(newChain2.Blocks[0].Users) != 1 {
		t.Fail()
	}

	//ensure that the new chain added the hashed credentials from the seed, not the explicit pass
	if newChain2.Blocks[0].Users[0] != "user:"+hashAuth("user", "pass") {
		t.Fail()
	}

	//blocks shouldn't have the same hashes (even if their timestamps are the same, which they are in this test func)
	if newChain2.Blocks[0].Hash == chain2.Blocks[0].Hash {
		t.Fail()
	}
}
