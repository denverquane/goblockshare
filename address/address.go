package address

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
)

type Base64Address string

func HashPublicToB64Address(pub ecdsa.PublicKey) Base64Address {
	h := sha256.New()
	h.Write(pub.X.Bytes())
	h.Write(pub.Y.Bytes())
	return Base64Address(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}