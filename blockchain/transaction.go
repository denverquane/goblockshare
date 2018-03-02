package blockchain

type Transaction struct {
	Author 		string
	Channel 	string
	Message 	string
}

func (trans Transaction) ToString() string {
	return trans.Author + " posted \"" + trans.Message + "\" on the " + trans.Channel + " channel"
}

func SampleTransaction() Transaction {
	return Transaction{"John Doe", "Test", "Sample message."}
}
