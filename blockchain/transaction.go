package blockchain

type Transaction struct {
	Author 		string
	Channel 	string
	Message 	string
}

func (trans Transaction) ToString() string {
	return trans.Author + " posted \"" + trans.Message + "\" on " + trans.Channel + " channel"
}
