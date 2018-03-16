//Blockchain contains the crucial aspects of maintaining a blockchain, such as Blocks, Transactions, and the operations
//these structures need and support
package blockchain

import (
	"encoding/hex"
	"crypto/sha256"
	"time"
	"strconv"
	"fmt"
	"log"
	"strings"
)

//Block represents the building "block" of the chain; any time a block is generated, it represents a change in the
//overall state of the chain, and successive blocks of the chain. For example, a user leaving a comment on a channel
//should be reflected in a new (immutable) block that other users would not be able to edit or remove; only "tack onto"
type Block struct {
	Index        int64
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
	Difficulty   int
	Nonce		 string
}

//ToString simply returns a human-legible representation of a Block in question
func (block Block) ToString() string {
	str := "Block: \n[\n   Index: " + strconv.Itoa(int(block.Index)) + "\n   Time: " + block.Timestamp +
		"\n   Total Transactions: " + strconv.Itoa(len(block.Transactions)) + "\n"
		str += "   Hash: " + block.Hash[0:5] + "...\n   PrevHash: "
		if len(block.PrevHash) > 4 {
			str += block.PrevHash[0:5] + "...\n]\n"
		} else {
			str += "...\n]\n"
		}
		return str
}

//InitialBlock creates a Block that has index 0, present timestamp, empty transaction slice,
//and an accurate/valid hash (albeit no previous hash for obvious reasons)
func InitialBlock() Block {
	var initBlock Block
	t := time.Now()
	initBlock.Index = 0
	initBlock.Timestamp = t.Format(time.RFC1123)
	initBlock.Transactions = make([]Transaction, 0)
	initBlock.PrevHash = ""
	initBlock.Hash = t.String()
	initBlock.Difficulty = 1

	initBlock.Hash = calcHash(initBlock)
	return initBlock
}

//calcHash calculates the hash for a given block based on ALL its attributes
func calcHash(block Block) string {
	record := string(block.Index) + block.Timestamp
	for _, v := range block.Transactions {
		record += v.ToString()
	}
	record += block.PrevHash + string(block.Difficulty) + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

//GenerateBlock accepts a "base" block to append to, and a transaction. The function
//creates a new block from the base block, and appends the transaction to it (rehashing and updating
//as necessary)
func GenerateBlock(oldBlock Block, transaction Transaction) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.Format(time.RFC1123)
	newBlock.Transactions = append(oldBlock.Transactions, transaction)
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = oldBlock.Difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		hash := calcHash(newBlock)

		if !isHashValid(hash, newBlock.Difficulty) {
			fmt.Println("Do more work: " + hash)
			continue
		} else {
			newBlock.Hash = hash
			break
		}
	}
	return newBlock
}

//IsBlockSequenceValid checks if an old block and a new block are capable of following one another;
//whether they form an valid chain of blocks or not. If the indexes or hashes (previous and current)
//do not match, then the sequence is invalid
func IsBlockSequenceValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		fmt.Println(newBlock.ToString() + "doesn't have the correct index to follow:\n" + oldBlock.ToString())
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		log.Println(newBlock.ToString() + "has a prevHash that doesn't match the hash of:" + oldBlock.ToString())
		return false
	}

	str := calcHash(newBlock)
	if str != newBlock.Hash {
		log.Println(newBlock.ToString() + "has a hash that doesn't match: " + str)
		return false
	}

	return true
}