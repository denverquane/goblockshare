package blockchain

import (
	"errors"
	"fmt"
	"log"
)

var globalChains = make(map[string]*BlockChain)

//CreateNewChannel attempts to create a new blockchain from a received transaction, and the name of the ADMIN/
//authorizing channel. If the transaction is invalid, or the user is not an authorized admin, then the function
//returns an empty blockchain and an error
func CreateNewGlobalChannel(transaction AuthTransaction) (BlockChain, error) {

	if transaction.TransactionType != "CREATE_CHANNEL" {
		return BlockChain{}, errors.New("Incorrect transtype for creating a channel")
	}

	if transaction.Message == "" {
		return BlockChain{}, errors.New("transaction message/channel name is empty")
	}

	if val, ok := globalChains[transaction.Message]; ok == false {
		chain := MakeInitialChain([]UserPassPair{{transaction.Username, transaction.Password}})
		globalChains[transaction.Message] = &chain
		log.Println("Created \"" + transaction.Message + "\" channel with initial admin: \"" + transaction.Username + "\"")
		return *globalChains[transaction.Message], nil
	} else {
		log.Println("Tried to create a channel that already exists")
		return *val, errors.New("Tried to create a channel that already exists")
	}
}

//SetChannelChain attempts to map the name of a channel to the blockchain that it represents. If there is already an
//entry for a specified channel name, the function returns the chain and an error (this function does not allow directly
//overriding an existing channel by the same name)
func SetChannelChain(channel string, chain BlockChain) (BlockChain, error) {
	if channel == "" {
		log.Println("Invalid empty channel name provided")
		return BlockChain{}, errors.New("Invalid empty channel name provided")
	}

	if len(chain.Blocks) == 0 {
		log.Println("Empty chain provided")
		return BlockChain{}, errors.New("Empty chain provided")
	}

	if val, ok := globalChains[channel]; ok == false {
		if chain.IsValid() {
			globalChains[channel] = &chain
			log.Println("Added \"" + channel + "\" channel")
			return *globalChains[channel], nil
		} else {
			log.Println("Attempted to set " + channel + " channel, but chain is invalid")
			return chain, errors.New("Attempted to set " + channel + " channel, but chain is invalid")
		}
	} else {
		log.Println("Tried to set the globalchain when one is already set...")
		return *val, errors.New("Tried to set the globalchain when one is already set...")
	}
}

//GetChainByValue returns the channel/blockchain associated with a provided string name
func GetChainByValue(channel string) (BlockChain, error) {
	if channel == "" {
		log.Println("Attempted to access channel with empty name")
		return BlockChain{}, errors.New("Channel name is empty")
	}
	if _, ok := globalChains[channel]; ok == true {
		return *globalChains[channel], nil
	} else {
		log.Printf("Attempted to access non-existent \"%s\" channel", channel)
		return BlockChain{}, errors.New(channel + " does not exist")
	}
}

//AttemptReplaceChain takes a channel name and a blockchain, and checks to see if the existing chain for the provided
//string name can be replaced with the newly provided blockchain. This function checks to ensure that the new chain is
//valid, of the same branch as the original channel, and is longer than the existing chain. If any of these conditions
//fail, the function returns the provided chain and an error
func AttemptReplaceChain(channel string, newChain BlockChain) (BlockChain, error) {
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
				return thisChain, errors.New("chains are of different branches or versions, keeping mine!")
			}
		} else {
			return thisChain, errors.New("provided chain is not longer than the current chain")
		}
	} else {
		return thisChain, errors.New("provided chain is invalid; keeping old chain")
	}
}

//WriteTransaction attempts to write a transaction to the channel referred to by the provided string name.
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

func GetChannelNames() []string {
	var retArr []string

	for i, _ := range globalChains {
		retArr = append(retArr, i)
	}
	return retArr
}
