package blockchain

import (
	"fmt"
	"github.com/pkg/errors"
)

var globalChain* BlockChain = nil
var Peers []string = nil

func SetGlobalChain(chain BlockChain) {
	if globalChain == nil {
		globalChain = &chain
	} else {
		fmt.Println("Tried to set the globalchain when one is already set...")
	}
}

func GetChainByValue() BlockChain {
	return *globalChain
}

func CheckReplacementChain(newChain BlockChain) (BlockChain, error) {
	var thisChain = GetChainByValue()

	if newChain.IsValid() {

		if newChain.Len() > thisChain.Len() {
			if AreChainsSameBranch(thisChain, newChain) {
				globalChain.mux.Lock()

				globalChain.AppendMissingBlocks(newChain)

				globalChain.mux.Unlock()

				return *globalChain, nil
			} else {
				return thisChain, errors.New("Chains are of different branches, keeping mine!")
			}
		} else {
			return thisChain, errors.New("Provided chain is not longer than the current chain")
		}
	} else {
		return thisChain, errors.New("Provided chain is invalid; keeping old chain")
	}
}

func WriteTransaction(trans AuthTransaction) (BlockChain, error) {
	var thisChain = GetChainByValue()

	if trans.TransactionType == "" || trans.Message == "" {
		return thisChain, errors.New("Supply transaction in this format: " + GetTransactionFormat())
	}

	oldBlock := thisChain.GetNewestBlock()
	newBlock, err := GenerateBlock(oldBlock, []AuthTransaction{trans})

	if err != nil {
		return thisChain, err
	}

	fmt.Println("New block:\n" + newBlock.ToString())

	if IsBlockSequenceValid(newBlock, oldBlock) {
		globalChain.mux.Lock()
		globalChain.Blocks = append(globalChain.Blocks, newBlock)
		globalChain.mux.Unlock()
		fmt.Println(globalChain.GenerateCollapsedChannelChat())
		return *globalChain, nil
		//Block = blockchain.CheckLongerChain(newBlock, Block)
		//fmt.Println("Successfully added: {" + m.RemovePassword().ToString() + "} to the chain")
		//BroadcastToAllPeers(Peers, *globalChain)
	} else {
		return thisChain, errors.New("Block sequence invalid somehow...")
	}
}





