package blockchain

import (
	"fmt"
	"testing"
)

func TestMakeInitialChain(t *testing.T) {
	chain := MakeInitialChain()

	fmt.Println(chain.ToString())

	if len(chain.Blocks) != 1 {
		t.Fail()
	}

	if len(chain.Blocks[0].Transactions) != 0 {
		t.Fail()
	}

	if chain.Blocks[0].Index != 0 {
		t.Fail()
	}

}
