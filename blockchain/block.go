package blockchain

import (
	"encoding/hex"
	"crypto/sha256"
	"time"
)

type Block struct {
	index 		int64
	timestamp 	string
	transactions []Transaction
	hash 		string
	prevHash 	string
}

func calcHash(block Block) string {
	record := string(block.index) + block.timestamp
	for _, v := range block.transactions {
		record += toString(v)
	}
	record += block.hash + block.prevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func checkLongerChain(block1, block2 Block) Block {
	if len(block1.transactions) > len(block2.transactions) {
		return block1
	}
	return block2
}

func generateBlock(oldBlock Block, transaction Transaction) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.index = oldBlock.index + 1
	newBlock.timestamp = t.String()
	newBlock.transactions = append(oldBlock.transactions, transaction)
	newBlock.prevHash = oldBlock.hash
	newBlock.hash = calcHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.index+1 != newBlock.index {
		return false
	}

	if oldBlock.hash != newBlock.prevHash {
		return false
	}

	if calcHash(newBlock) != newBlock.hash {
		return false
	}

	return true
}
