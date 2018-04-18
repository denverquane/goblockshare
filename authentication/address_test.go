package authentication

import (
	"testing"
	"fmt"
	trans "github.com/denverquane/GoBlockShare/transaction"
)

func TestStuff(t *testing.T) {
	userdets := GenerateNewPersonalAddress()
	//fmt.Println("Private: " + userdets.private.D.String() + "\n" + "Public: " + userdets.public.Y.String())

	simpleT := trans.MakeSimple("sample message", trans.ADD_MESSAGE)

	fmt.Println(userdets.address)

	fullTrans := trans.MakeFull(simpleT, userdets.publicKey, userdets.privateKey, "XPvMLIArcc48ctx26EDwgRdtt72mXW6bqEhqY6xiFT8=")
	fmt.Println(fullTrans.ToString())
}