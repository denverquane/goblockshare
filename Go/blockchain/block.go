//Blockchain contains the crucial aspects of maintaining a blockchain, such as Blocks, Transactions, and the operations
//these structures need and support
package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//Block represents the building "block" of the chain; any time a block is generated, it represents a change in the
//overall state of the chain, and successive blocks of the chain. For example, a user leaving a comment on a channel
//should be reflected in a new (immutable) block that other users would not be able to edit or remove; only "tack onto"
type Block struct {
	Index        int64
	Timestamp    string
	Transactions []Transaction
	Users        []string
	Hash         string
	PrevHash     string
	Difficulty   int
	Nonce        string
}

type UserPassPair struct {
	Username string
	Password string
}

//ToString simply returns a human-legible representation of a Block in question
func (block Block) ToString() string {
	str := "Block: \n[\n   Index: " + strconv.Itoa(int(block.Index)) + "\n   Time: " + block.Timestamp +
		"\n   Total Transactions: " + strconv.Itoa(len(block.Transactions)) + "\n" + "   Users: \n"
	for _, v := range block.Users {
		str += "     " + v + "\n"
	}
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
func InitialBlock(users []UserPassPair) Block {
	var initBlock Block
	t := time.Now()
	initBlock.Index = 0
	initBlock.Timestamp = t.Format(time.RFC1123)
	initBlock.Transactions = make([]Transaction, 0)
	initBlock.Users = make([]string, len(users))
	for i, v := range users {
		initBlock.Users[i] = v.Username + ":" + hashAuth(v.Username, v.Password)
	}
	//initBlock.PrevHash = "GoBlockShare Version: " + version
	initBlock.Hash = t.String() //placeholder until we calculate the actual hash
	initBlock.Difficulty = 1

	initBlock = initBlock.hashUntilValid(6)

	return initBlock
}

func hashAuth(username, password string) string {
	h := sha256.New()
	h.Write([]byte(username + password))
	return hex.EncodeToString(h.Sum(nil))
}

//hashUntilValid continually increments a block's "Nonce" until the block hashes correctly to the provided
//difficulty
func (block Block) hashUntilValid(difficulty int) Block {
	block.Hash = calcHash(block)

	for i := 0; !isHashValid(block.Hash, difficulty); i++ {
		hexx := fmt.Sprintf("%x", i)
		block.Nonce = hexx
		block.Hash = calcHash(block)
	}
	return block
}

//calcHash calculates the hash for a given block based on ALL its attributes
func calcHash(block Block) string {
	record := string(block.Index) + block.Timestamp
	for _, v := range block.Transactions {
		record += v.ToString()
	}
	for _, v := range block.Users {
		record += v
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
func GenerateBlock(oldBlock Block, transactions []AuthTransaction) (Block, error) {

	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.Format(time.RFC1123)
	newBlock.Users = oldBlock.Users
	newBlock.Difficulty = oldBlock.Difficulty
	newBlock.PrevHash = oldBlock.Hash

	for _, t := range transactions {
		if !t.IsValidType() {
			log.Println("Invalid transaction type supplied!!!")
			// log.Println(t.ToString())
			fmt.Println("Retaining old block")
			return oldBlock, errors.New("Invalid type supplied")
		}

		if !t.IsAuthorized(oldBlock.Users) {
			log.Println("User is not authorized!!!")
			// log.Println(t.ToString())
			fmt.Println("Retaining old block")
			return oldBlock, errors.New("User not authorized")
		}

		if t.TransactionType == "ADD_USER" {
			cleanTrans, err := t.VerifyAndFormatAddUserTrans(oldBlock)

			if err != nil {
				return oldBlock, err
			}

			newBlock.Users = append(newBlock.Users, cleanTrans.Message)
			newBlock.Transactions = append(newBlock.Transactions, cleanTrans)
			newBlock.Difficulty += 1 // The larger a user list becomes, the harder it should become to post messages
		} else {
			newBlock.Transactions = append(newBlock.Transactions, t.RemovePassword())
		}
	}

	hashTime := time.Now()
	newBlock = newBlock.hashUntilValid(newBlock.Difficulty)
	endTime := time.Now()
	fmt.Println("Took " + strconv.Itoa(endTime.Second()-hashTime.Second()) + " seconds to mine with diff=" +
		strconv.Itoa(newBlock.Difficulty*len(transactions)))

	return newBlock, nil
}

//ValidateAddUser ensures that the message provided not only matches the correct format for adding a user,
//but also doesn't contain an entry for adding a user that already exists
func (oldBlock Block) ValidateAddUser(message string) (string, error) {
	strs := strings.Split(message, ":")
	if len(strs) < 2 {
		return "", errors.New("Parse error of user/pass in string: " + message)
	}

	user := strs[0]
	pass := strs[1]

	for _, v := range oldBlock.Users {
		u := strings.Split(v, ":")[0]
		if u == user {
			return "", errors.New("User \"" + user + "\" is already registered!")
		}
	}

	return user + ":" + hashAuth(user, pass), nil
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
