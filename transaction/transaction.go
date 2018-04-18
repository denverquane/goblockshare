package transaction

import (
	"math/big"
	"crypto/ecdsa"
	"crypto/rand"
	"log"
	"github.com/denverquane/GoBlockShare/address"
)

type TransType int

const (
	TEST TransType = iota
	ADD_MESSAGE
	ADD_FILE
	DELETE_MESSAGE
	ADD_USER
	CREATE_CHANNEL
)

/*
type Transaction interface {
	OriginPubKey() ecdsa.PublicKey
	DestinationAddr() address.Base64Address
}
*/

type SimpleTransaction struct {
	message string
	tType TransType
}

type FullTransaction struct {
	originPubKey ecdsa.PublicKey
	originAddr address.Base64Address
	txRef	[]string
	signedPayload	SignedTransaction
	destination address.Base64Address
}

func MakeFull(s SimpleTransaction, key ecdsa.PublicKey, priv ecdsa.PrivateKey, dest address.Base64Address) FullTransaction{
	signed := s.SignMessage(&priv)
	return FullTransaction{key, address.HashPublicToB64Address(key), []string{}, signed, dest}
}

func (ft FullTransaction) Verify() bool {
	if !ft.signedPayload.VerifyWithKey(ft.originPubKey) { //signed transaction isn't verified with the public key
		return false
	} else if address.HashPublicToB64Address(ft.originPubKey) != ft.originAddr { //public key does not match up with the address
		return false
	}
	// TODO Right here, verify the history of transactions to the origin address
	return true
}

func (ft FullTransaction) ToString() string {
	return "Public key " + ft.originPubKey.X.String() + ft.originPubKey.Y.String() + "\nand address: " + string(ft.originAddr) +
		"\n sending " + ft.signedPayload.simple.message + "\nto " + string(ft.destination)
}

func MakeSimple(message string, tt TransType) SimpleTransaction {
	return SimpleTransaction{message, tt}
}

func (s SimpleTransaction) GetMessageBytes() []byte {
	return []byte(s.message)
}

func (h SimpleTransaction) SignMessage(priv *ecdsa.PrivateKey) SignedTransaction {
	r, s, err := ecdsa.Sign(rand.Reader, priv, h.GetMessageBytes())

	if err != nil {
		log.Println("Error when signing transaction!")
		return SignedTransaction{}
	}
	return SignedTransaction{h, r, s}
}


type SignedTransaction struct {
	simple SimpleTransaction
	r, s *big.Int
}

func (s SignedTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	return ecdsa.Verify(&key, s.simple.GetMessageBytes(), s.r, s.s)
}

type UnconfirmedTransaction struct {
	Signed SignedTransaction
	PublicKey ecdsa.PublicKey
	Address address.Base64Address
}

func (uc UnconfirmedTransaction) Verify() bool {
	if !uc.Signed.VerifyWithKey(uc.PublicKey) { //signed transaction isn't verified with the public key
		return false
	} else if address.HashPublicToB64Address(uc.PublicKey) != uc.Address { //public key does not match up with the address
		return false
	}
	// TODO Right here, verify the history of transactions to the origin address
	return true
}

