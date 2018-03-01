package blockchain

type Transaction struct {
	author 		string
	channel 	string
	message 	string
}

func toString(trans Transaction) string {
	return trans.author + " posted \"" + trans.message + "\" on " + trans.message + " channel"
}
