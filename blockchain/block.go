package blockchain

import (
	"encoding/hex"
	"crypto/sha256"
	"time"
	"strconv"
	"fmt"
)

type Block struct {
	Index        int64
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
}

type BlockChain struct {
	Length		int
	Blocks		[]Block
}

func InitialBlock() Block {
	var initBlock Block
	t := time.Now()
	initBlock.Index = 0
	initBlock.Timestamp = t.String()
	initBlock.Transactions = make([]Transaction, 0)
	initBlock.PrevHash = ""
	initBlock.Hash = t.String()

	initBlock.Hash = calcHash(initBlock)
	return initBlock
}

func calcHash(block Block) string {
	record := string(block.Index) + block.Timestamp
	for _, v := range block.Transactions {
		record += v.ToString()
	}
	record += block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (block Block) ToString() string {
	var str string
	str =  "Block @ index " + strconv.Itoa(int(block.Index)) + " was created at " + block.Timestamp +
		" and has transactions: \n[\n";

		for _, v := range block.Transactions {
			str += v.ToString() + "\n"
		}
		str += "]\nWith Hash: " + block.Hash + " and prevhash: " + block.PrevHash

		return str
}

func GenerateBlock(oldBlock Block, transaction Transaction) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transactions = append(oldBlock.Transactions, transaction)
	newBlock.PrevHash = oldBlock.Hash

	newBlock.Hash = calcHash(newBlock)

	return newBlock, nil
}


func IsBlockSequenceValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	str := calcHash(newBlock)
	if str != newBlock.Hash {
		fmt.Println("Mismatch hash!")
		fmt.Println("New: " + str + " Old: " + newBlock.Hash)

		return false
	}

	return true
}

func (chain BlockChain) IsValid() bool {
	if chain.Length != len(chain.Blocks) {
		return false
	}

	if chain.Length < 2 {
		return true
	}

	for i := 0; i < chain.Length-1; i++ {
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
	if chain1.Length > chain2.Length {
		min = chain2.Length
	} else {
		min = chain1.Length
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
	return chain.Blocks[chain.Length-1]
}
