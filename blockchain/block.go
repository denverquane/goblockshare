//Blockchain contains the crucial aspects of maintaining a blockchain, such as Blocks, Transactions, and the operations
//these structures need and support
package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Block represents the building "block" of the chain; any time a block is generated, it represents a change in the
//overall state of the chain, and successive blocks of the chain. For example, a user leaving a comment on a channel
//should be reflected in a new (immutable) block that other users would not be able to edit or remove; only "tack onto"
type Block struct {
	Index        int64
	Timestamp    string
	Transactions []transaction.FullTransaction
	Hashed         string
	PrevHash     string
	Difficulty   int
	Nonce        string
	mux          sync.Mutex
}

//ToString simply returns a human-legible representation of a Block in question
func (block Block) ToString() string {
	str := "Block: \n[\n   Index: " + strconv.Itoa(int(block.Index)) + "\n   Time: " + block.Timestamp +
		"\n   Total Transactions: " + strconv.Itoa(len(block.Transactions)) + "\n"
	str += "   Hash: " + block.Hashed[0:5] + "...\n   PrevHash: "
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
	initBlock.Transactions = make([]transaction.FullTransaction, 0)
	//initBlock.PrevHash = "GoBlockShare Version: " + version
	initBlock.Hashed = t.String() //placeholder until we calculate the actual hash
	initBlock.Difficulty = 1

	for i := 0; !isHashValid(initBlock.Hashed, 3); i++ {
		hexx := fmt.Sprintf("%x", i)
		initBlock.Nonce = hexx
		initBlock.Hashed = initBlock.Hash()
	}

	return initBlock
}

//hashUntilValid continually increments a block's "Nonce" until the block hashes correctly to the provided
//difficulty
func (block *Block) hashUntilValid(difficulty int, c chan bool) {
	block.mux.Lock()
	block.Hashed = block.Hash()
	block.mux.Unlock()

	for i := 0; !isHashValid(block.Hashed, difficulty); i++ {
		c <- false
		hexx := fmt.Sprintf("%x", i)
		block.mux.Lock()
		block.Nonce = hexx
		block.Hashed = block.Hash()
		block.mux.Unlock()
	}
	c <- true
}

func (block *Block) AddTransaction(trans transaction.FullTransaction) {
	fmt.Println("Adding transaction to mining block")
	block.mux.Lock()
	block.Transactions = append(block.Transactions, trans)
	block.mux.Unlock()
}

//calcHash calculates the hash for a given block based on ALL its attributes
func (block Block) Hash() string {

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

//GenerateBlock expects a "base" block to append transactions to, and thus "mining" a new block that contains these
//transactions. The difficulty in mining this new block is proportional to the number of transactions being added,
//and the more users that are registered to a channel results in a higher difficulty for mining transactions
func GenerateInvalidBlock(oldBlock Block, transactions []transaction.FullTransaction) (Block, error) {

	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.Format(time.RFC1123)
	newBlock.Difficulty = oldBlock.Difficulty
	newBlock.PrevHash = oldBlock.Hashed

	for _, t := range transactions {
		if !t.Verify() {
			log.Println("Invalid transaction!!!")
			log.Println(t.ToString())
			fmt.Println("Retaining old block")
			return oldBlock, errors.New("Invalid transaction supplied")
		}

		newBlock.Transactions = append(newBlock.Transactions, t)
	}

	return newBlock, nil
}

//IsBlockSequenceValid checks if an old block and a new block are capable of following one another;
//whether they form an valid chain of blocks or not. If the indexes or hashes (previous and current)
//do not match, then the sequence is invalid
func IsBlockSequenceValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		fmt.Println(newBlock.ToString() + "doesn't have the correct index to follow:\n" + oldBlock.ToString())
		return false
	}

	if oldBlock.Hashed != newBlock.PrevHash {
		log.Println(newBlock.ToString() + "has a prevHash that doesn't match the hash of:" + oldBlock.ToString())
		return false
	}

	str := newBlock.Hash()
	if str != newBlock.Hashed {
		log.Println(newBlock.ToString() + "has a hash that doesn't match: " + str)
		return false
	}

	return true
}
