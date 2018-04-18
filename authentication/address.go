package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"
	"strconv"
	addr "github.com/denverquane/GoBlockShare/address"
)

var AUTHENTICATION_CURVE = elliptic.P256()

// See https://golang.org/src/crypto/ecdsa/ecdsa_test.go

type PersonalAddress struct {
	privateKey ecdsa.PrivateKey
	publicKey ecdsa.PublicKey
	address addr.Base64Address
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
	address := addr.HashPublicToB64Address(priv.PublicKey)

	if strings.Contains(string(address), "O") || strings.Contains(string(address), "l"){
		fmt.Println("invalid char, regenerating")
		return GenerateNewPersonalAddress()
	}

	fmt.Println("New Key Pair created")
	fmt.Println("Encoded length: " + strconv.Itoa(len(address)))
	return PersonalAddress{*priv, priv.PublicKey, addr.HashPublicToB64Address(priv.PublicKey)}
}