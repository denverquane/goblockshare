package main

import (
	"github.com/denverquane/GoBlockShare/address"
	"fmt"
)

func main() {
	genericTesting()
}

func genericTesting() {
	addr := address.GenerateNewPersonalAddress()
	fmt.Println(addr.GetB64Address())
}
