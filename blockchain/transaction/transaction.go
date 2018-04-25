package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
)

type Transaction interface {
	ToString() string
}

type OriginInfo struct {
	PubKeyX big.Int
	PubKeyY big.Int
	Address Base64Address
}

func AddressToOriginInfo(address PersonalAddress) OriginInfo {
	return OriginInfo{*address.PublicKey.X, *address.PublicKey.Y, address.Address}
}

type RESTWrappedFullTransaction struct {
	Origin   OriginInfo
	Txref    []string
	Quantity float64
	Payload  string
	R        big.Int
	S        big.Int
	DestAddr string
}

func (oi OriginInfo) GetHash() []byte {
	h := sha256.New()
	h.Write(oi.PubKeyX.Bytes())
	h.Write(oi.PubKeyY.Bytes())
	h.Write([]byte(oi.Address))
	return h.Sum(nil)
}

type SignedTransaction struct {
	Origin   OriginInfo
	DestAddr Base64Address
	Quantity float64
	Payload  string
	R, S     *big.Int
}

func (st SignedTransaction) GetHash(haveRSbeenSet bool) []byte {
	h := sha256.New()
	h.Write(st.Origin.GetHash())
	h.Write([]byte(st.DestAddr))
	// -1 as the precision arg gets the # to 64bit precision intuitively
	h.Write([]byte(strconv.FormatFloat(st.Quantity, 'f', -1, 64)))

	//Filters the cases where we just want the hash for non-signing purposes
	//(if the transaction hasn't been signed, we shouldn't hash R and S as they don't matter)
	if haveRSbeenSet {
		h.Write(st.R.Bytes())
		h.Write(st.S.Bytes())
	}
	h.Write([]byte(st.Payload))
	return h.Sum(nil)
}

type FullTransaction struct {
	SignedTrans SignedTransaction
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

func (rest RESTWrappedFullTransaction) ConvertToFull() (FullTransaction, error) {
	var full = FullTransaction{}
	full.SignedTrans = SignedTransaction{rest.Origin, Base64Address(rest.DestAddr), rest.Quantity, rest.Payload, &rest.R, &rest.S}
	full.TxRef = rest.Txref
	full.TxID = hex.EncodeToString(full.GetHash())
	return full, nil
}

func (st SignedTransaction) MakeFull(txref []string) FullTransaction {
	full := FullTransaction{st, txref, ""}
	full.TxID = hex.EncodeToString(full.GetHash())
	return full
}

func (st SignedTransaction) SignMessage(priv *ecdsa.PrivateKey) SignedTransaction {
	hashed := st.GetHash(false)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)

	if err != nil {
		log.Println("Error when signing transaction!")
		return SignedTransaction{}
	}
	st.R = r
	st.S = s
	return st
}

func (s SignedTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	return ecdsa.Verify(&key, s.GetHash(false), s.R, s.S)
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
	return true
}

func (oi OriginInfo) ToString() string {
	return "\n{\n\"origin\":\n{\n\"address\":\"" + string(oi.Address) + "\",\n\"pubkeyx\":" + oi.PubKeyX.String() +
		",\n\"pubkeyy\":" + oi.PubKeyY.String() + "\n},\n"
}

func (st SignedTransaction) ToString() string {
	return st.Origin.ToString() + "\"txref\":[],\n\"quantity\":" +
		strconv.FormatFloat(st.Quantity, 'f', -1, 64) + ",\n\"payload\":\"" +
		st.Payload + "\",\n\"r\":" + st.R.String() + ",\n\"s\":" + st.S.String() + ",\n\"destAddr\":\"" +
		string(st.DestAddr) + "\"\n}\n"
}
