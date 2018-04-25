package transaction

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
)

type Transaction interface {
	ToString() string
}

type OriginInfo struct {
	OriginPubKeyX big.Int
	OriginPubKeyY big.Int
	OrigAddr string
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
//func SignMessage(orig Base64Address, dest Base64Address, quantity float64, str string, priv *ecdsa.PrivateKey) SignedTransaction {
//	r, s, err := ecdsa.Sign(rand.Reader, priv, []byte(str))
//
//	if err != nil {
//		log.Println("Error when signing transaction!")
//		return SignedTransaction{}
//	}
//	return SignedTransaction{orig, dest, quantity, str,  r, s}
//}

func (s SignedTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	return ecdsa.Verify(&key, []byte(s.payload), s.R, s.S)
}

func (st SignedTransaction) Verify() bool {
	key := ecdsa.PublicKey{AUTHENTICATION_CURVE, &st.Origin.OriginPubKeyX, &st.Origin.OriginPubKeyY}

	if !st.VerifyWithKey(key) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(key) != Base64Address(st.Origin.OrigAddr) { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	// TODO Right here, verify the history of transactions to the origin address
	return true
}

//TODO This needs to be fast when hashing many transactions into a single block
//func (ft FullTransaction) ToString() string {
//	return "Public key " + ft.OriginPubKeyX.String() + ft.OriginPubKeyY.String() + "\nand address: " + string(ft.OriginAddr) +
//		"\n sending " + ft.SignedPayload.payload + "\nto " + string(ft.Destination)
//}

func (oi OriginInfo) ToString() string {
	return oi.OrigAddr + " with x=" + oi.OriginPubKeyX.String() + " and y=" + oi.OriginPubKeyY.String()
}

func (st SignedTransaction) ToString() string {
	return st.Origin.ToString()
}

func (oi OriginInfo) Hash() []byte {
	h := sha256.New()
	h.Write(oi.OriginPubKeyX.Bytes())
	h.Write(oi.OriginPubKeyY.Bytes())
	h.Write([]byte(oi.OrigAddr))
	return h.Sum(nil)
}

func (st SignedTransaction) Hash() []byte {
	h := sha256.New()
	h.Write(st.Origin.Hash())
	h.Write([]byte(st.destAddr))
	//TODO fix these
	//h.Write([]byte(st.quantity))
	//h.Write([]byte(st.R))
	//h.Write([]byte(st.S))
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
