package lib

import "testing"

func TestJWTToken(t *testing.T) {
	username := "user"
	token, err := CreateToken(username)
	if err != nil {
		t.Error("Could not create a token, err: ", err.Error())
		return
	}

	usernameFromToken, err := ValidateToken(token)
	if err != nil {
		t.Error("Could not validate token, err: ", err.Error())
		return
	}

	if usernameFromToken != username {
		t.Error("The username from the jwt token doesn't match the correct username.")
		return
	}
}
