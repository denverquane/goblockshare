package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
)

// This struct is for storing the information required to interact with a channel
type ChannelRecord struct {
	haveToken         bool
	isFullyConfigured bool

	channelPublic  crypto.PublicKey
	channelPrivate crypto.PrivateKey
	channelAddress transaction.Base64Address

	myPublic  crypto.PublicKey
	myPrivate crypto.PrivateKey
	myAddress transaction.Base64Address
}

func GenerateNewChannelRecord(tokenName string, address transaction.Base64Address, chain blockchain.BlockChain, amt float64) ChannelRecord {
	record := ChannelRecord{haveToken: false, isFullyConfigured: false}

	private, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Println("Problem generating RSA key!")
		return record
	}

	record.myPrivate = private
	record.myPublic = private.Public()
	record.myAddress = address

	record.haveToken = amt > 0.0
	fmt.Println("Generated a new channel record for: " + tokenName)

	return record
}
