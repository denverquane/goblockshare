package blockchain

import (

)

type TransType int

const (
	ADD_MESSAGE TransType = iota
	DELETE_MESSAGE
	ADD_USER
)

var ValidTransactionTypes = map[TransType]string{
	ADD_MESSAGE : "ADD_MESSAGE",
	DELETE_MESSAGE : "DELETE_MESSAGE",
	ADD_USER: "ADD_USER",
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
	TransactionType string
}

type UserTransaction struct {
	Username		string
	Message 		string
}

func (trans AuthTransaction) IsValidType() bool {
	for _,v := range ValidTransactionTypes {
		if trans.TransactionType == v {
			return true
		}
	}
	return false
}

func (trans AuthTransaction) IsAuthorized(authUsers []string) bool {
	var auth = trans.Username + ":" + hashAuth(trans.Username, trans.Password)

	for _, v := range authUsers {
		if auth == v {
			return true
		}
	}

	return false
}

func (trans AuthTransaction) CensorAddUserTrans(userHash string) Transaction {
	return Transaction{Username:trans.Username, Channel:"",
		Message:userHash, TransactionType:trans.TransactionType}
}

func (trans AuthTransaction) RemovePassword() Transaction {
	return Transaction{Username:trans.Username, Channel:trans.Channel,
		Message:trans.Message, TransactionType:trans.TransactionType}
}

func (trans Transaction) ToString() string {
	return trans.Username + " posted \"" + trans.Message + "\" on the " + trans.Channel + " channel"
}

func SampleAuthTransaction(user, pass string) AuthTransaction {
	return AuthTransaction{user, pass, "Test", "Sample message.", "ADD_MESSAGE"}
}

func GetTransactionFormat() string {
	return "{username:User,password:Pass,channel:TestChannel,message:SampleMessage,transactiontype:ADD_MESSAGE}"
}
