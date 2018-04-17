package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

// See https://golang.org/src/crypto/ecdsa/ecdsa_test.go

type UserDetails struct {
	private ecdsa.PrivateKey
	public ecdsa.PublicKey
}

func GenerateNewUserDetails() UserDetails {
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)

	if err != nil {

		fmt.Errorf("error: %s", err)

		return UserDetails{}

	}

	if !c.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		fmt.Errorf("public key invalid: %s", err)
	}
	return UserDetails{*priv, priv.PublicKey}

}
