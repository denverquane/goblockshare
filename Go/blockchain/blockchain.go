package blockchain

import (
	"sync"
)

type BlockChain struct {
	Blocks []Block
	mux    sync.Mutex
}

func (chain BlockChain) Len() int {
	return len(chain.Blocks)
}

func (chain BlockChain) ToString() string {
	str := "Chain: \n{\n"
	for _, v := range chain.Blocks {
		str += v.ToString()
	}
	str += "}"
	return str
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

// AreChainsNewAndSameBranch verifies that the chains are both new chains, and at least have the same users
func AreChainsNewAndSameBranch(chain1, chain2 BlockChain) bool {
	if chain1.Len() == 1 || chain2.Len() == 1 {
		if calcHash(chain1.Blocks[0]) != chain1.Blocks[0].Hash ||
			calcHash(chain2.Blocks[0]) != chain2.Blocks[0].Hash {
				return false
		}

		if len(chain1.Blocks[0].Users) != len(chain2.Blocks[0].Users) {
			return false
		} else {
			for i, v := range chain1.Blocks[0].Users {
				if v != chain2.Blocks[0].Users[i] {
					return false // different user lists
				}
			}
			return true // same lists of initial users
		}
	} else {
		return false
	}
}

func (chain BlockChain) GetNewestBlock() Block {
	return chain.Blocks[chain.Len()-1]
}

func MakeInitialChain(users []UserPassPair) BlockChain {
	chain := BlockChain{Blocks: make([]Block, 1)}
	chain.Blocks[0] = InitialBlock(users)
	return chain
}

//AppendMissingBlocks takes a chain, and appends all the transactions that are found on a longer chain to it
//This is handy when using a single Global chain that should never be entirely replaced; only appended to
func (chain BlockChain) AppendMissingBlocks(longerChain BlockChain) BlockChain {
	if AreChainsSameBranch(chain, longerChain) && longerChain.IsValid() {
		for i := len(chain.Blocks); i < len(longerChain.Blocks); i++ {
			chain.Blocks = append(chain.Blocks, longerChain.Blocks[i])
		}
	}
	return chain
}
