package blockchain

type Block struct {
	index 		int64
	timestamp 	int64
	transactions []Transaction
	proof 		int64
	prev_hash 	string
}


