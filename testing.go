package main

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
)

func main() {
	genericTesting()
}

func genericTesting() {
	addr := transaction.GenerateNewPersonalAddress()
	fmt.Println("X", addr.PublicKey.X)
	fmt.Println("Y", addr.PublicKey.Y)
	s := transaction.SimpleTransaction{"dsfgsd", transaction.ADD_MESSAGE}
	signed := s.SignMessage(&addr.PrivateKey)
	fmt.Println(signed.Simple.Message)
	fmt.Println("r", signed.R)
	fmt.Println("s", signed.S)
	fmt.Println(addr.Address)
}
