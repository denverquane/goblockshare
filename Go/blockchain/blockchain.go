package blockchain

import (
	"sync"
)

type BlockChain struct {
	Blocks		[]Block
	mux 		sync.Mutex
}

var ChannelTransMap = make(map[string][]UserTransaction)
var IndexCounter = -1

func (chain BlockChain) Len() int {
	return len(chain.Blocks)
}

func (chain BlockChain) IsValid() bool {
	if chain.Len() != len(chain.Blocks) {
		return false
	}

	if chain.Len() < 2 {
		return true
	}

	for i := 0; i < chain.Len()-1; i++ {
		oldB := chain.Blocks[i]
		newB := chain.Blocks[i+1]

		if !IsBlockSequenceValid(newB, oldB) {
			return false
		}
	}
	return true
}

func AreChainsSameBranch(chain1, chain2 BlockChain) bool {
	var min = 0
	if chain1.Len() > chain2.Len() {
		min = chain2.Len()
	} else {
		min = chain1.Len()
	}
	for i := 0; i < min; i++ {
		a := chain1.Blocks[i]
		b := chain2.Blocks[i]
		if calcHash(a) != calcHash(b) {
			return false
		}
	}
	return true
}

func (chain BlockChain) GetNewestBlock() Block {
	return chain.Blocks[chain.Len()-1]
}

func MakeInitialChain(users []UserPassPair, version string) BlockChain {
	chain := BlockChain{Blocks: make([]Block, 1)}
	chain.Blocks[0] = InitialBlock(users, version)
	return chain
}

//AppendMissingBlocks takes a chain, and appends all the transactions that are found on a longer chain to it
//This is handy when using a single Global chain that should never be entirely replaced; only appended to
func (chain BlockChain) AppendMissingBlocks (longerChain BlockChain) {
	if AreChainsSameBranch(chain, longerChain) && longerChain.IsValid(){
		for i := len(chain.Blocks); i < len(longerChain.Blocks); i++ {
			chain.Blocks = append(chain.Blocks, longerChain.Blocks[i])
		}
	}
}

func (chain BlockChain) GenerateCollapsedChannelChat() map[string][]UserTransaction{
	block := chain.GetNewestBlock()
	if !(block.Index > int64(IndexCounter)) {
		return ChannelTransMap //blocks haven't been updated, just return the old one
	}

	//TODO Update index, and only process the transactions not already processed in the map
	for _, v := range block.Transactions {
		if v.TransactionType != ADD_USER {
			if _, ok := ChannelTransMap[v.Channel]; ok { // Is the entry (the channel name) found in the map?
				ChannelTransMap[v.Channel] = append(ChannelTransMap[v.Channel], UserTransaction{v.Username, v.Message})
				//if it's found, add the message to that channel
				//if it's found, add the message to that channel
			} else { //otherwise, make the channel in the map, and then add the transaction
				ChannelTransMap[v.Channel] = make([]UserTransaction, 1)
				ChannelTransMap[v.Channel][0] = UserTransaction{v.Username, v.Message}
			}
		}
	}
	return ChannelTransMap
}