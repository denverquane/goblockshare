package main

import "github.com/denverquane/GoBlockShare/blockchain/transaction"

type Wallet struct {
	addresses []transaction.PersonalAddress
	balanceMap map[string]float64
}



