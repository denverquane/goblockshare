package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
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

//TODO Alter signed so it includes the origin pubkey, origin addr, and destination!
//Needed to verify all these important aspects of the msg for security, but also for simplicity
type SignedTransaction struct {
	Simple string
	R, S   *big.Int
}

type FullTransaction struct {
	OriginPubKeyX  *big.Int
	OriginPubKeyY  *big.Int
	OriginAddr    Base64Address
	TxRef         []string
	SignedPayload SignedTransaction
	Destination   Base64Address
}

//TODO test with sending sub-objects? Would simplify the format of the REST API/conversions considerably
type RESTWrappedFullTransaction struct {
	OriginPubKeyX string
	OriginPubKeyY string
	OriginAddress string
	Txref         []string
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
	full.OriginPubKeyX = x
	full.OriginPubKeyY = y
	full.OriginAddr = Base64Address(rest.OriginAddress)
	full.TxRef = rest.Txref
	r := new(big.Int)
	r.SetString(rest.R, 10)
	s := new(big.Int)
	s.SetString(rest.S, 10)
	full.SignedPayload = SignedTransaction{rest.SignedMsg, r, s}
	full.Destination = Base64Address(rest.DestAddr)
	return full
}

//func MakeFull(str string, origin PersonalAddress, dest Base64Address) (FullTransaction, error) {
//	signed := SignMessage(str, &origin.PrivateKey)
//	full := FullTransaction{origin.PublicKey, origin.Address, []string{}, signed, dest}
//	if !full.Verify() {
//		return FullTransaction{}, errors.New("Generated transaction is invalid!")
//	}
//	return full, nil
//}

func SignMessage(str string, priv *ecdsa.PrivateKey) SignedTransaction {
	r, s, err := ecdsa.Sign(rand.Reader, priv, []byte(str))

	if err != nil {
		log.Println("Error when signing transaction!")
		return SignedTransaction{}
	}
	return SignedTransaction{str, r, s}
}

func (s SignedTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	return ecdsa.Verify(&key, []byte(s.Simple), s.R, s.S)
}

func (s SignedTransaction) ToString() string {
	return s.Simple
}

func (ft FullTransaction) Verify() bool {
	key := ecdsa.PublicKey{AUTHENTICATION_CURVE, ft.OriginPubKeyX, ft.OriginPubKeyY}

	if !ft.SignedPayload.VerifyWithKey(key) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(key) != ft.OriginAddr { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	// TODO Right here, verify the history of transactions to the origin address
	return true
}

//TODO This needs to be fast when hashing many transactions into a single block
func (ft FullTransaction) ToString() string {
	return "Public key " + ft.OriginPubKeyX.String() + ft.OriginPubKeyY.String() + "\nand address: " + string(ft.OriginAddr) +
		"\n sending " + ft.SignedPayload.Simple + "\nto " + string(ft.Destination)
}
