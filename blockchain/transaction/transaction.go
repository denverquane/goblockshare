package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

type SignableTransaction interface {
	GetHash(bool) []byte
	SetRS(*big.Int, *big.Int) SignableTransaction
	GetRS() (*big.Int, *big.Int)
	GetOrigin() OriginInfo
	ToString() string
}

func Sign(priv *ecdsa.PrivateKey, st SignableTransaction) SignableTransaction {
	hashed := st.GetHash(false)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)

	if err != nil {
		log.Println("Error when signing transaction!")
		return st
	}
	return st.SetRS(r, s)
}

func Verify(st SignableTransaction) bool {
	origin := st.GetOrigin()
	key := ecdsa.PublicKey{AUTHENTICATION_CURVE, &origin.PubKeyX, &origin.PubKeyY}

	if !VerifyWithKey(st, key) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(key) != Base64Address(origin.Address) { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	return true
}

func VerifyWithKey(st SignableTransaction, key ecdsa.PublicKey) bool {
	r, s := st.GetRS()
	return ecdsa.Verify(&key, st.GetHash(false), r, s)
}

func MakeFull(st SignableTransaction, txref []string) FullTransaction {
	full := FullTransaction{st, txref, ""}
	full.TxID = hex.EncodeToString(full.GetHash())
	return full
}

type OriginInfo struct {
	PubKeyX big.Int
	PubKeyY big.Int
	Address Base64Address
}

func (oi OriginInfo) GetHash() []byte {
	h := sha256.New()
	h.Write(oi.PubKeyX.Bytes())
	h.Write(oi.PubKeyY.Bytes())
	h.Write([]byte(oi.Address))
	return h.Sum(nil)
}

func (oi OriginInfo) ToString() string {
	return "\n{\n\"origin\":\n{\n\"address\":\"" + string(oi.Address) + "\",\n\"pubkeyx\":" + oi.PubKeyX.String() +
		",\n\"pubkeyy\":" + oi.PubKeyY.String() + "\n},\n"
}

func AddressToOriginInfo(address PersonalAddress) OriginInfo {
	return OriginInfo{*address.PublicKey.X, *address.PublicKey.Y, address.Address}
}

type RESTWrappedFullTransaction struct {
	Origin   OriginInfo
	Txref    []string
	Quantity float64
	Currency string
	Payload  string
	R        big.Int
	S        big.Int
	DestAddr string
}

func (rest RESTWrappedFullTransaction) ConvertToFull() (FullTransaction, error) {
	var signed = SignedTransaction{rest.Origin, Base64Address(rest.DestAddr), rest.Quantity, rest.Currency, rest.Payload, &rest.R, &rest.S}
	var full = MakeFull(signed, rest.Txref)
	return full, nil
}

type FullTransaction struct {
	SignedTrans SignableTransaction
	TxRef       []string
	TxID        string
}

func (ft FullTransaction) GetHash() []byte {
	h := sha256.New()
	h.Write(ft.SignedTrans.GetHash(true))
	for _, v := range ft.TxRef {
		h.Write([]byte(v))
	}
	return h.Sum(nil)
}
