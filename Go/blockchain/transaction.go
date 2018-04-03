package blockchain

import (
	"strings"
	"errors"
)

type TransType int

const (
	TEST TransType = iota
	ADD_MESSAGE
	ADD_FILE
	DELETE_MESSAGE
	ADD_USER
)

var ValidTransactionTypes = map[string]TransType{
	"TEST" : TEST,
	"ADD_MESSAGE" : ADD_MESSAGE,
	"ADD_FILE" : ADD_FILE,
	"DELETE_MESSAGE" : DELETE_MESSAGE,
	"ADD_USER" : ADD_USER,
}

type AuthTransaction struct {
	Username 		string
	Password		string
	Channel 		string
	Message 		string
	TransactionType	string
}

type Transaction struct {
	Username        string
	Channel         string
	Message         string
	TransactionType TransType
}

type UserTransaction struct {
	Username		string
	Message 		string
}

func (trans AuthTransaction) IsValidType() bool {
	for i,_ := range ValidTransactionTypes {
		if trans.TransactionType == i {
			return true
		}
	}
	return false
}

// IsAuthorized verifies that the transaction is being posted by a user who is found in the list
// of authorized users
func (trans AuthTransaction) IsAuthorized(authUsers []string) bool {
	var auth = trans.Username + ":" + hashAuth(trans.Username, trans.Password)

	for _, v := range authUsers {
		if auth == v {
			return true
		}
	}

	return false
}

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
	return Transaction{Username:trans.Username, Channel:trans.Channel,
		Message:trans.Message, TransactionType:ValidTransactionTypes[trans.TransactionType]}
}

func (trans Transaction) ToString() string {
	return trans.Username + " posted \"" + trans.Message + "\" on the " + trans.Channel + " channel"
}

func (trans AuthTransaction) ToString() string {
	return trans.Username + " w/ pass: " + trans.Password + " posted " + trans.TransactionType + " type to " +
		trans.Channel + " w/ message: " + trans.Message
}

func SampleAuthTransaction(user, pass string) AuthTransaction {
	return AuthTransaction{user, pass, "Test", "Sample message.", "TEST"}
}

func GetTransactionFormat() string {
	return "{username:User,password:Pass,channel:TestChannel,message:SampleMessage,transactiontype:ADD_MESSAGE}"
}
