package tests

import (
	"goutils/data"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestCreateToken(t *testing.T) {
	secret := "secretkey"
	datastring := "{ \"username\" : \"yourusername\" }"

	token, err := data.CreateToken(secret, datastring)
	if err != nil {
		t.Fatal(err)
	}

	if len(token) == 0 {
		t.Fatal("failed to create token")
	}

}
