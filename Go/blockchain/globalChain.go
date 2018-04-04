package blockchain

import (
	"errors"
	"fmt"
	"log"
)

var globalChains = make(map[string]*BlockChain)

func CreateNewChannel(transaction AuthTransaction, adminChainName string) (BlockChain, error) {
	adminChain, err := GetChainByValue(adminChainName)
	if err != nil {
		return BlockChain{}, err //don't ever return the admin chain until authorization
	}
	if !transaction.IsAuthorized(adminChain.GetNewestBlock().Users) {
		return BlockChain{}, errors.New("Not authorized to post to admin channel")
	}

	if transaction.TransactionType != "CREATE_CHANNEL" {
		return BlockChain{}, errors.New("Incorrect transtype for creating a channel")
	}

	if val, ok := globalChains[transaction.Message]; ok == false {
		chain := CreateChainFromSeed(adminChain)
		globalChains[transaction.Message] = &chain
		log.Println("Created " + transaction.Message + " channel")
		return *globalChains[transaction.Message], nil
	} else {
		log.Println("Tried to create a channel that already exists")
		return *val, errors.New("Tried to create a channel that already exists")
	}
}

func SetChannelChain(channel string, chain BlockChain) (BlockChain, error) {
	if val, ok := globalChains[channel]; ok == false {
		globalChains[channel] = &chain
		log.Println("Added \"" + channel + "\" channel")
		return *globalChains[channel], nil
	} else {
		log.Println("Tried to set the globalchain when one is already set...")
		return *val, errors.New("Tried to set the globalchain when one is already set...")
	}
}

func GetChainByValue(channel string) (BlockChain, error) {
	if _, ok := globalChains[channel]; ok == true {
		return *globalChains[channel], nil
	} else {
		log.Printf("Attempted to access non-existent \"%s\" channel", channel)
		return BlockChain{}, errors.New(channel + " does not exist")
	}
}

func CheckReplacementChain(channel string, newChain BlockChain) (BlockChain, error) {
	var thisChain, err = GetChainByValue(channel)

	if err != nil {
		return thisChain, err
	}

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
	var thisChain, err = GetChainByValue(channel)

	if err != nil {
		return thisChain, err
	}

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
