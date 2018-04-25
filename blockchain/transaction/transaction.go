package transaction

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
	"log"
)

type Transaction interface {
	ToString() string
}

type OriginInfo struct {
	PubKeyX big.Int
	PubKeyY big.Int
	Address string
}

type RESTWrappedFullTransaction struct {
	Origin   	  OriginInfo
	Txref         []string
	Quantity      float64
	Payload		  string
	R             big.Int
	S             big.Int
	DestAddr      string
}

type SignedTransaction struct {
	Origin 	 OriginInfo
	destAddr Base64Address
	quantity float64
	payload  string
	R, S   	 *big.Int
}

type FullTransaction struct {
	SignedTrans SignedTransaction
	TxRef         []string
	TxID		  string
}

func (rest RESTWrappedFullTransaction) ConvertToFull() (FullTransaction, error) {
	var full = FullTransaction{}
	full.SignedTrans = SignedTransaction{rest.Origin, Base64Address(rest.DestAddr), rest.Quantity, rest.Payload, &rest.R, &rest.S}
	full.TxRef = rest.Txref
	full.TxID = hex.EncodeToString(full.Hash())
	return full, nil
}

//func MakeFull(str string, origin PersonalAddress, dest Base64Address) (FullTransaction, error) {
//	signed := SignMessage(str, &origin.PrivateKey)
//	full := FullTransaction{origin.PublicKey, origin.Address, []string{}, signed, dest}
//	if !full.Verify() {
//		return FullTransaction{}, errors.New("Generated transaction is invalid!")
//	}
//	return full, nil
//}

// TODO Sign every aspect of the transaction (sign an unsigned transaction? -> have a getBytes()?)
func (st SignedTransaction) SignMessage(priv *ecdsa.PrivateKey) SignedTransaction {
	var arr = []byte(st.Origin.Hash())
	arr = append(arr, []byte(st.destAddr)...)
	// arr = append(arr, byte[](strconv.FormatFloat(st.quantity, 64))...)
	arr = append(arr, []byte(st.payload)...)
	arr = append(arr, st.R.Bytes()...)
	arr = append(arr, st.S.Bytes()...)
	r, s, err := ecdsa.Sign(rand.Reader, priv, arr)

	if err != nil {
		log.Println("Error when signing transaction!")
		return SignedTransaction{}
	}
	st.R = r
	st.S = s
	return st
}

func (s SignedTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	return ecdsa.Verify(&key, []byte(s.payload), s.R, s.S)
}

func (st SignedTransaction) Verify() bool {
	key := ecdsa.PublicKey{AUTHENTICATION_CURVE, &st.Origin.PubKeyX, &st.Origin.PubKeyY}

	if !st.VerifyWithKey(key) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(key) != Base64Address(st.Origin.Address) { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	// TODO Right here, verify the history of transactions to the origin address
	return true
}

//TODO This needs to be fast when hashing many transactions into a single block
//func (ft FullTransaction) ToString() string {
//	return "Public key " + ft.OriginPubKeyX.String() + ft.PubKeyY.String() + "\nand address: " + string(ft.OriginAddr) +
//		"\n sending " + ft.SignedPayload.payload + "\nto " + string(ft.Destination)
//}

func (oi OriginInfo) ToString() string {
	return oi.Address + " with x=" + oi.PubKeyX.String() + " and y=" + oi.PubKeyY.String()
}

func (st SignedTransaction) ToString() string {
	return st.Origin.ToString()
}

func (oi OriginInfo) Hash() []byte {
	h := sha256.New()
	h.Write(oi.PubKeyX.Bytes())
	h.Write(oi.PubKeyY.Bytes())
	h.Write([]byte(oi.Address))
	return h.Sum(nil)
}

func (st SignedTransaction) Hash() []byte {
	h := sha256.New()
	h.Write(st.Origin.Hash())
	h.Write([]byte(st.destAddr))
	//TODO fix these
	//h.Write([]byte(st.quantity))
	h.Write(st.R.Bytes())
	h.Write(st.S.Bytes())
	h.Write([]byte(st.payload))
	return h.Sum(nil)
}

func (ft FullTransaction) Hash() []byte {
	h := sha256.New()
	h.Write(ft.SignedTrans.Hash())
	for _, v := range ft.TxRef {
		h.Write([]byte(v))
	}
	return h.Sum(nil)
}
