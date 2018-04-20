package transaction

import (
	"math/big"
	"crypto/ecdsa"
	"crypto/rand"
	"log"
	"github.com/denverquane/GoBlockShare/blockchain/transaction/address"
	"errors"
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

type Transaction interface {
	ToString() string
}

type SimpleTransaction struct {
	Message string
	TType TransType
}

type SignedTransaction struct {
	simple SimpleTransaction
	r, s *big.Int
}

type FullTransaction struct {
	originPubKey ecdsa.PublicKey
	originAddr address.Base64Address
	txRef	[]string
	signedPayload	SignedTransaction
	destination address.Base64Address
}

func MakeFull(s SimpleTransaction, origin address.PersonalAddress, dest address.Base64Address) (FullTransaction, error) {
	signed := s.SignMessage(&origin.PrivateKey)
	full := FullTransaction{origin.PublicKey, origin.Address, []string{}, signed, dest}
	if !full.Verify() {
		return FullTransaction{}, errors.New("Generated transaction is invalid!")
	}
	return full, nil
}

func (s SimpleTransaction) GetMessageBytes() []byte {
	return []byte(s.Message)
}

func (h SimpleTransaction) SignMessage(priv *ecdsa.PrivateKey) SignedTransaction {
	r, s, err := ecdsa.Sign(rand.Reader, priv, h.GetMessageBytes())

	if err != nil {
		log.Println("Error when signing transaction!")
		return SignedTransaction{}
	}
	return SignedTransaction{h, r, s}
}

func (s SimpleTransaction) ToString() string {
	return "Message: " + s.Message + " of type: " + string(s.TType)
}

func (s SignedTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	return ecdsa.Verify(&key, s.simple.GetMessageBytes(), s.r, s.s)
}

func (s SignedTransaction) ToString() string {
	return s.simple.ToString()
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
		"\n sending " + ft.signedPayload.simple.Message + "\nto " + string(ft.destination)
}

