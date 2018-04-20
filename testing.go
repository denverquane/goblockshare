package main

import (
	"github.com/denverquane/GoBlockShare/blockchain/transaction/address"
	"fmt"
)

func main() {
	genericTesting()
}

func genericTesting() {
	addr := address.GenerateNewPersonalAddress()
	fmt.Println(addr.Address)
}
