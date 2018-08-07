package blockchain

import (
	"github.com/denverquane/goblockshare/common"
	"testing"
)

func TestMakeInitialChain(t *testing.T) {
	chain := MakeInitialChain()

	if len(chain.Blocks) != 1 {
		t.Fail()
	}

	if chain.Blocks[0].Index != 0 {
		t.Fail()
	}

	if chain.IsProcessing() {
		t.Fail()
	}

	if !chain.IsValid() {
		t.Fail()
	}

	if chain.ToString() == "" {
		t.Fail()
	}
}

func TestBlockChain_GetTxById(t *testing.T) {
	chain := MakeInitialChain()

	if chain.GetTxById("sample").TxID != "ERROR" {
		t.Fail()
	}
}

//func TestAreChainsSameBranch(t *testing.T) {
//
//	f := func(chain BlockChain) BlockChain {
//		addr := common.GenerateNewPersonalAddress()
//		origin := addr.ConvertToOriginInfo()
//		torr := common.SetAliasTrans{"sf"}
//		signable := common.NewSignable(origin, torr, common.SET_ALIAS)
//		added, err := chain.AddTransaction(signable, "")
//		if added || err == nil {
//			t.Fail()
//		}
//		signed := signable.SignAndSetTxID(&addr.PrivateKey)
//		added, err = chain.AddTransaction(signed, "")
//		if !added || err != nil {
//			t.Fail()
//		}
//		return chain
//	}
//	chain1 := MakeInitialChain()
//	chain1 = f(chain1)
//
//	chain2 := chain1
//	chain2 = f(chain2)
//
//	if !AreChainsSameBranch(chain1, chain2) {
//		t.Fail()
//	}
//}

func TestBlockChain_GetNewestBlock(t *testing.T) {
	chain := MakeInitialChain()
	chain.GetNewestBlock()
}

func TestBlockChain_AddTransaction(t *testing.T) {
	chain := MakeInitialChain()
	addr := common.GenerateNewPersonalAddress()
	origin := addr.ConvertToOriginInfo()
	torr := common.SetAliasTrans{"sf"}
	signable := common.NewSignable(origin, torr, common.SET_ALIAS)
	added, err := chain.AddTransaction(signable, "")
	if added || err == nil {
		t.Fail()
	}
	signed := signable.SignAndSetTxID(&addr.PrivateKey)
	added, err = chain.AddTransaction(signed, "")
	if !added || err != nil {
		t.Fail()
	}
}
