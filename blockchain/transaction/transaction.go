package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
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
	TType   TransType
}

type SignedTransaction struct {
	Simple SimpleTransaction
	R, S   *big.Int
}

type FullTransaction struct {
	originPubKey  ecdsa.PublicKey
	originAddr    Base64Address
	txRef         []string
	signedPayload SignedTransaction
	destination   Base64Address
}

type RESTWrappedFullTransaction struct {
	OriginPubKeyX string
	OriginPubKeyY string
	OriginAddress string
	SignedMsg     string
	R             string
	S             string
	DestAddr      string
}

func (rest RESTWrappedFullTransaction) ConvertToFull() FullTransaction {
	var full = FullTransaction{}
	x := new(big.Int)
	x.SetString(rest.OriginPubKeyX, 10)
	y := new(big.Int)
	y.SetString(rest.OriginPubKeyY, 10)
	full.originPubKey = ecdsa.PublicKey{AUTHENTICATION_CURVE, x, y}
	full.originAddr = Base64Address(rest.OriginAddress)
	full.txRef = []string{}
	r := new(big.Int)
	r.SetString(rest.R, 10)
	s := new(big.Int)
	s.SetString(rest.S, 10)
	full.signedPayload = SignedTransaction{SimpleTransaction{rest.SignedMsg, ADD_MESSAGE}, r, s}
	full.destination = Base64Address(rest.DestAddr)
	return full
}

func MakeFull(s SimpleTransaction, origin PersonalAddress, dest Base64Address) (FullTransaction, error) {
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
	return ecdsa.Verify(&key, s.Simple.GetMessageBytes(), s.R, s.S)
}

func (s SignedTransaction) ToString() string {
	return s.Simple.ToString()
}

func (ft FullTransaction) Verify() bool {
	if !ft.signedPayload.VerifyWithKey(ft.originPubKey) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(ft.originPubKey) != ft.originAddr { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	// TODO Right here, verify the history of transactions to the origin address
	return true
}

func (ft FullTransaction) ToString() string {
	return "Public key " + ft.originPubKey.X.String() + ft.originPubKey.Y.String() + "\nand address: " + string(ft.originAddr) +
		"\n sending " + ft.signedPayload.Simple.Message + "\nto " + string(ft.destination)
}
