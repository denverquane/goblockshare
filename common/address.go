package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"math/big"
	"log"
)

type Base64Address string

var AUTHENTICATION_CURVE = elliptic.P256()

// See https://golang.org/src/crypto/ecdsa/ecdsa_test.go

type PersonalAddress struct {
	PrivateKey ecdsa.PrivateKey //TODO Encrypt and save locally
	PublicKey  ecdsa.PublicKey
	Address    Base64Address
}

type OriginInfo struct {
	PubKeyX *big.Int
	PubKeyY *big.Int
	Address Base64Address
}

func (oi OriginInfo) GetRawBytes() []byte {
	return []byte(string(oi.PubKeyX.Bytes()) + string(oi.PubKeyY.Bytes()) + string(oi.Address))
}

func (oi OriginInfo) ToString() string {
	return "\n{\n\"Origin\":\n{\n\"Address\":\"" + string(oi.Address) + "\",\n\"Pubkeyx\":" + oi.PubKeyX.String() +
		",\n\"Pubkeyy\":" + oi.PubKeyY.String() + "\n},\n"
}

func AddressToOriginInfo(address PersonalAddress) OriginInfo {
	return OriginInfo{address.PublicKey.X, address.PublicKey.Y, address.Address}
}

func GenerateNewPersonalAddress() PersonalAddress {
	priv, err := ecdsa.GenerateKey(AUTHENTICATION_CURVE, rand.Reader)

	if err != nil {

		log.Fatal(fmt.Errorf("error: %s", err))

		return PersonalAddress{}

	}

	if !AUTHENTICATION_CURVE.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		log.Fatal(fmt.Errorf("public key invalid: %s", err))
	}
	address := HashPublicToB64Address(priv.PublicKey)

	//get rid of "O" and "l" letters to avoid ambiguity with numbers
	if strings.Contains(string(address), "O") || strings.Contains(string(address), "l") || strings.Contains(string(address), "/"){
		//fmt.Println("invalid char, regenerating")
		return GenerateNewPersonalAddress()
	}

	/******************* Testing encoding a public/private key pair in pem format *********************/
	//PubASN1, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	// do something about it
	//}
	//
	//pubBytes := pem.EncodeToMemory(&pem.Block{
	//	Type:  "ECDSA PUBLIC KEY",
	//	Bytes: PubASN1,
	//})
	//fmt.Println(string(pubBytes))
	/**************************************************************************************************/

	//fmt.Println("New Key Pair created")
	return PersonalAddress{*priv, priv.PublicKey, HashPublicToB64Address(priv.PublicKey)}
}

func HashPublicToB64Address(pub ecdsa.PublicKey) Base64Address {
	h := sha256.New()
	h.Write(pub.X.Bytes())
	h.Write(pub.Y.Bytes())
	return Base64Address(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}
