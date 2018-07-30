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

func (js JSONSignableTransaction) ConvertToSignable() SignableTransaction {
	copied := SignableTransaction{}

	copied.Origin = js.Origin
	copied.TransactionType = js.TransactionType
	copied.TxID = js.TxID
	copied.R = js.R
	copied.S = js.S

	switch js.TransactionType {
	case "PUBLISH_TORRENT":
		var mm PublishTorrentTrans
		if err := json.Unmarshal([]byte(js.Transaction), &mm); err != nil {
			log.Fatal(err)
		}
		copied.Transaction = mm
		break
	case "TORRENT_REP":
		var mm TorrentRepTrans
		if err := json.Unmarshal([]byte(js.Transaction), &mm); err != nil {
			log.Fatal(err)
		}
		copied.Transaction = mm
		break
	case "SHARED_LAYER":
		var mm SharedLayerTrans
		if err := json.Unmarshal([]byte(js.Transaction), &mm); err != nil {
			log.Fatal(err)
		}
		copied.Transaction = mm
		break
	case "LAYER_REP":
		var mm LayerRepTrans
		if err := json.Unmarshal([]byte(js.Transaction), &mm); err != nil {
			log.Fatal(err)
		}
		copied.Transaction = mm
		break
	}
	return copied
}

func (st SignableTransaction) setRS(r *big.Int, s *big.Int) SignableTransaction {
	st.R = r
	st.S = s
	return st
}

func (st SignableTransaction) GetRS() (*big.Int, *big.Int) {
	return st.R, st.S
}

func (st SignableTransaction) GetHash() []byte {
	h := sha256.New()

	st.TxID = "" //don't hash the TXid; the txid is just the hash of everything else anyways
	h.Write([]byte(st.ToString()))

	return h.Sum(nil)
}

func (st SignableTransaction) GetOrigin() OriginInfo {
	return st.Origin
}

func (st SignableTransaction) ToString() string {
	data, _ := json.Marshal(st)
	return hex.EncodeToString(data)
}

func (st SignableTransaction) SignAndSetTxID(priv *ecdsa.PrivateKey) SignableTransaction {
	hashed := st.GetHash()
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)

	if err != nil {
		log.Println("Error when signing transaction!")
		return st
	}
	st = st.setRS(r, s)
	st.TxID = hex.EncodeToString(st.GetHash())
	return st
}

func (st SignableTransaction) Verify() bool {
	id := hex.EncodeToString(st.GetHash())
	if id != st.TxID {
		fmt.Println("Transaction doesn't hash to it's ID!")
		return false
	}

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
	return ecdsa.Verify(&key, st.GetHash(), r, s)
}

