package authentication

import (
	"testing"
	"fmt"
)

func TestGenerateNewUserDetails(t *testing.T) {
	userdets := GenerateNewUserDetails()
	fmt.Println("Private: " + userdets.private.D.String() + "\n" + "Public: " + userdets.public.Y.String())
}