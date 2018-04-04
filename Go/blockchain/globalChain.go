package blockchain

import (
	"errors"
	"fmt"
	"log"
)

var globalChains = make(map[string]*BlockChain)

func SetChannelChain(channel string, chain BlockChain) {
	if _, ok := globalChains[channel]; ok == false {
		globalChains[channel] = &chain
		log.Println("Added \"" + channel + "\" channel")
	} else {
		fmt.Println("Tried to set the globalchain when one is already set...")
	}
}

func GetChainByValue(channel string) BlockChain {
	return *globalChains[channel]
}

func CheckReplacementChain(channel string, newChain BlockChain) (BlockChain, error) {
	var thisChain = GetChainByValue(channel)

	if newChain.IsValid() {

		if newChain.Len() > thisChain.Len() {
			if AreChainsSameBranch(thisChain, newChain) {
				globalChains[channel].mux.Lock()

				globalChains[channel].AppendMissingBlocks(newChain)

				globalChains[channel].mux.Unlock()

				return *globalChains[channel], nil
			} else {
				return thisChain, errors.New("Chains are of different branches, keeping mine!")
			}
		} else {
			return thisChain, errors.New("provided chain is not longer than the current chain")
		}
	} else {
		return thisChain, errors.New("Provided chain is invalid; keeping old chain")
	}
}

func WriteTransaction(channel string, trans AuthTransaction) (BlockChain, error) {
	var thisChain = GetChainByValue(channel)

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
		globalChains[channel].mux.Lock()
		globalChains[channel].Blocks = append(globalChains[channel].Blocks, newBlock)
		globalChains[channel].mux.Unlock()
		return *globalChains[channel], nil
		//Block = blockchain.CheckLongerChain(newBlock, Block)
		//fmt.Println("Successfully added: {" + m.RemovePassword().ToString() + "} to the chain")
		//BroadcastToAllPeers(Peers, *globalChain)
	} else {
		return thisChain, errors.New("Block sequence invalid somehow...")
	}
}
