package blockchain

type TransType int

var ValidTransactionTypes = []string{
	"ADD_MESSAGE",
	"DELETE_MESSAGE",
	"ADD_USER"}


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


func (trans AuthTransaction) RemovePassword() Transaction {
	return Transaction{Username:trans.Username, Channel:trans.Channel,
		Message:trans.Message, TransactionType:trans.TransactionType}
}

func (trans Transaction) ToString() string {
	return trans.Username + " posted \"" + trans.Message + "\" on the " + trans.Channel + " channel"
}

func SampleAuthTransaction() AuthTransaction {
	return AuthTransaction{"user", "pass", "Test", "Sample message.", "ADD_MESSAGE"}
}

func GetTransactionFormat() string {
	return "{username:User,password:Pass,channel:TestChannel,message:SampleMessage,transactiontype:ADD_MESSAGE}"
}
