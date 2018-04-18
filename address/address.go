package address

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

type Base64Address string

var AUTHENTICATION_CURVE = elliptic.P256()

// See https://golang.org/src/crypto/ecdsa/ecdsa_test.go

type PersonalAddress struct {
	privateKey ecdsa.PrivateKey
	publicKey ecdsa.PublicKey
	address Base64Address
}

func GenerateNewPersonalAddress() PersonalAddress {
	priv, err := ecdsa.GenerateKey(AUTHENTICATION_CURVE, rand.Reader)

	if err != nil {

		fmt.Errorf("error: %s", err)

		return PersonalAddress{}

	}

	if !AUTHENTICATION_CURVE.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		fmt.Errorf("public key invalid: %s", err)
	}
	address := hashPublicToB64Address(priv.PublicKey)

	if strings.Contains(string(address), "O") || strings.Contains(string(address), "l"){
		fmt.Println("invalid char, regenerating")
		return GenerateNewPersonalAddress()
	}

	fmt.Println("New Key Pair created")
	return PersonalAddress{*priv, priv.PublicKey, hashPublicToB64Address(priv.PublicKey)}
}

func (p PersonalAddress) GetB64Address() Base64Address {
	return hashPublicToB64Address(p.publicKey)
}

func hashPublicToB64Address(pub ecdsa.PublicKey) Base64Address {
	h := sha256.New()
	h.Write(pub.X.Bytes())
	h.Write(pub.Y.Bytes())
	return Base64Address(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}