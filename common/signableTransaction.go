package common

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type SignableTransaction struct {
	Origin      OriginInfo // needed to say who I am (WITHIN the transaction)
	Transaction TorrentTransaction
	TransactionType string
	R, S        *big.Int // signature of the transaction, should be separate from the actual "message" components
	TxID        string
}

type JSONSignableTransaction struct {
	Origin      OriginInfo // needed to say who I am (WITHIN the transaction)
	Transaction json.RawMessage
	TransactionType string
	R, S        *big.Int // signature of the transaction, should be separate from the actual "message" components
	TxID        string
}

func (st SignableTransaction) setRS(r *big.Int, s *big.Int) SignableTransaction {
	st.R = r
	st.S = s
	return st
}

func (st SignableTransaction) GetRS() (*big.Int, *big.Int) {
	return st.R, st.S
}

func (st SignableTransaction) GetHash(haveRSbeenSet bool) []byte {
	h := sha256.New()
	h.Write(st.Origin.GetRawBytes())

	//if the transaction is merely a request for a resource (like a layer), we only really need the signature, not any
	//sort of nested transaction
	if st.Transaction != nil {
		h.Write(st.Transaction.GetRawBytes())
	}
	h.Write([]byte(st.TransactionType))


	//Filters the cases where we just want the hash for non-signing purposes
	//(if the transaction hasn't been signed, we shouldn't hash R and S as they don't matter)
	if haveRSbeenSet{
		h.Write(st.R.Bytes())
		h.Write(st.S.Bytes())
	}

	//don't hash the TXid; the txid is just the hash of everything else anyways
	return h.Sum(nil)
}

func (st SignableTransaction) GetOrigin() OriginInfo {
	return st.Origin
}

func (st SignableTransaction) ToString() string {
	return st.Origin.ToString() + "\n\"TransactionType\": \"" + st.TransactionType + "\",\n\"Transaction\":\n{\n" +
		string(st.Transaction.ToString()) + "\n},\n\"R\":" + st.R.String() + ",\n\"S\":" +
		st.S.String() + ",\n\"TxID\":\"" + st.TxID +"\"\n}\n"
}

func (st SignableTransaction) SignAndSetTxID(priv *ecdsa.PrivateKey) SignableTransaction {
	hashed := st.GetHash(false)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)

	if err != nil {
		log.Println("Error when signing transaction!")
		return st
	}
	st = st.setRS(r, s)
	st.TxID = hex.EncodeToString(st.GetHash(true))
	return st
}

func (st SignableTransaction) Verify() bool {
	origin := st.GetOrigin()
	key := ecdsa.PublicKey{AUTHENTICATION_CURVE, origin.PubKeyX, origin.PubKeyY}

	if st.VerifyWithKey(key) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(key) != Base64Address(origin.Address) { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	return true
}

func (st SignableTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	r, s := st.GetRS()
	return ecdsa.Verify(&key, st.GetHash(true), r, s)
}

