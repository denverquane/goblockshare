package blockchain

import (
	"strings"
	"testing"
)

func TestSampleAuthTransaction(t *testing.T) {
	trans := SampleAuthTransaction("", "")

	//sample trans should always be a valid type
	if !trans.IsValidType() {
		t.Fail()
	}

	//should never be authorized if the auth list is empty
	if trans.IsAuthorized([]string{}) {
		t.Fail()
	}
}

func TestAuthTransaction_RemovePassword(t *testing.T) {
	aTrans := SampleAuthTransaction("user", "pass")

	trans := aTrans.RemovePassword()

	if trans.Username != aTrans.Username {
		t.Fail()
	}

	if trans.TransactionType != ValidTransactionTypes[aTrans.TransactionType] {
		t.Fail()
	}

	if trans.Message != aTrans.Message {
		t.Fail()
	}
}

func TestAuthTransaction_CensorAddUserMessage(t *testing.T) {
	aTrans := AuthTransaction{Username: "user", Password: "pass", Message: "newuser:password", TransactionType: "ADD_USER"}
	block := InitialBlock([]UserPassPair{})

	stripped, err := aTrans.VerifyAndFormatAddUserTrans(block)

	if err != nil {
		t.Fail()
	}

	if strings.Contains(stripped.Message, "password") {
		t.Fail()
	}

}
