package blockchain

import (
	"errors"
	"strings"
)

type TransType int

const (
	TEST TransType = iota
	ADD_MESSAGE
	ADD_FILE
	DELETE_MESSAGE
	ADD_USER
	CREATE_CHANNEL
)

var ValidTransactionTypes = map[string]TransType{
	"TEST":           TEST,
	"ADD_MESSAGE":    ADD_MESSAGE,
	"ADD_FILE":       ADD_FILE,
	"DELETE_MESSAGE": DELETE_MESSAGE,
	"ADD_USER":       ADD_USER,
	"CREATE_CHANNEL": CREATE_CHANNEL,
}

// AuthTransaction represents a transaction posted with credentials supplied, such as from a web UI or POST request.
// This sort of transaction is identical to the more general "Transaction", except that the transactiontype is a string
// (so that it can be supplied externally), and the Password is expected to be hashed/removed before posted to the chain
type AuthTransaction struct {
	Username        string
	Password        string
	Message         string
	TransactionType string
}

// Transaction represents a more generic version of the "AuthTransaction", and is a datatype intended for internal
// processing, and posting to the blockchain (it doesn't contain a password that could be compromised by posting to the
// chain).
type Transaction struct {
	Username        string
	Message         string
	TransactionType TransType
}

type UserTransaction struct {
	Username string
	Message  string
}

// Ensures that the transaction is one of the registered transaction types
func (trans AuthTransaction) IsValidType() bool {
	for i, _ := range ValidTransactionTypes {
		if trans.TransactionType == i {
			return true
		}
	}
	return false
}

// IsAuthorized verifies that the transaction is being posted by a user who is found in the list
// of authorized users for a channel
func (trans AuthTransaction) IsAuthorized(authUsers []string) bool {
	if trans.Username == "" || trans.Password == "" {
		return false
	}

	var auth = trans.Username + ":" + hashAuth(trans.Username, trans.Password)

	for _, v := range authUsers {
		if auth == v {
			return true
		}
	}

	return false
}

// VerifyAndFormatAddUserTrans verifies, upon a transaction to add a user to a channel/blockchain, that the transaction
// has the correct message formatting and that the user is not already found in the channel's user listings. Additionally,
// this function returns a "Transaction" (with proper message formatting) to guarantee that the transaction does not
// compromise the password
func (trans AuthTransaction) VerifyAndFormatAddUserTrans(oldBlock Block) (Transaction, error) {
	strs := strings.Split(trans.Message, ":")
	if len(strs) < 2 {
		return Transaction{}, errors.New("Parse error of user/pass in string: " + trans.Message)
	}

	user := strs[0]
	pass := strs[1]

	for _, v := range oldBlock.Users {
		u := strings.Split(v, ":")[0]
		if u == user {
			return Transaction{}, errors.New("User \"" + user + "\" is already registered!")
		}
	}

	strippedTrans := trans.RemovePassword()
	strippedTrans.Message = user + ":" + hashAuth(user, pass)
	return strippedTrans, nil
}

// RemovePassword takes an authorized transaction and converts it to a standard transaction,
// so it can be securely posted to the blockchain without posting a plaintext authorization
func (trans AuthTransaction) RemovePassword() Transaction {
	return Transaction{Username: trans.Username, Message: trans.Message,
		TransactionType: ValidTransactionTypes[trans.TransactionType]}
}

func (trans Transaction) ToString() string {
	return trans.Username + " posted \"" + trans.Message + "\""
}

func (trans AuthTransaction) ToString() string {
	return trans.Username + " w/ pass: " + trans.Password + " posted " + trans.TransactionType +
		" w/ message: " + trans.Message
}

func SampleAuthTransaction(user, pass string) AuthTransaction {
	return AuthTransaction{user, pass, "Sample message.", "TEST"}
}

func GetTransactionFormat() string {
	return "{username:User,password:Pass,message:SampleMessage,transactiontype:ADD_MESSAGE}"
}
