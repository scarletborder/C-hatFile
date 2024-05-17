package auth_test

import (
	"chatFileBackend/handlers/auth"
	"testing"
)

func TestLoginVerify(t *testing.T) {
	timestamp := 1715872756148
	enc2 := "71a120228ae52c477db27cc0dd1bd20a77b882b243f035d0c5c675e3eecc654b"
	// password := "!!8964jss"

	if !auth.LoginVerify("xjp", enc2, timestamp) {
		t.Error("not equal, expected", enc2)
	}
}
