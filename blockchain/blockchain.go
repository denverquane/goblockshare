package blockchain

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
)

type BlockChain struct {
	Blocks          []Block
	processingBlock *Block
}

func (chain BlockChain) IsProcessing() bool {
	return chain.processingBlock != nil
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

func (chain *BlockChain) AddTransaction(trans transaction.FullTransaction) {
	if chain.processingBlock != nil { //currently processing a block
		chain.processingBlock.AddTransaction(trans)
		return
	} else {
		invalidBlock, err := GenerateInvalidBlock(chain.GetNewestBlock(), []transaction.FullTransaction{trans})
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		var c = make(chan bool)
		chain.processingBlock = &invalidBlock
		fmt.Println("Mining a new block")
		go chain.processingBlock.hashUntilValid(5, c)
		for i := 0; !(<-c); i++ {
			if i%100000 == 0 {
				fmt.Println("Mining...")
			}
			// Wait until block is mined successfully
		}
		fmt.Println("Successfully mined block!")

		chain.Blocks = append(chain.Blocks, *chain.processingBlock)
		// fmt.Println(len(chain.Blocks[1].Transactions))
		chain.processingBlock = nil
	}
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

func MakeInitialChain() BlockChain {
	chain := BlockChain{Blocks: make([]Block, 1)}
	chain.Blocks[0] = InitialBlock()
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
