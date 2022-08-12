package auth

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken("test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token: ", token)
}