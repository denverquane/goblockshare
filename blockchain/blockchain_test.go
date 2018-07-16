package blockchain

import (
	"fmt"
	"testing"
)

func TestMakeInitialChain(t *testing.T) {
	chain := MakeInitialChain()

	//fmt.Println(chain.ToString())

	if len(chain.Blocks) != 1 {
		fmt.Println("No blocks!")
		t.Fail()
	}

	//TODO disabled for preliminary testing of Token payouts
	//if len(chain.Blocks[0].Transactions) != 1 {
	//	fmt.Println("There isn't exactly 1 transaction for the initial block")
	//	t.Fail()
	//}

	if chain.Blocks[0].Index != 0 {
		fmt.Println("Index isn't 0!")
		t.Fail()
	}

}
