package lib

import "testing"

func TestPasswordHash(t *testing.T) {
	startingPassword := "secret"
	hashed, err := HashPassword(startingPassword)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if hashed == "" {
		t.Error("The hashed password is empty.")
		return
	}

	if !CheckPasswordHash(startingPassword, hashed) {
		t.Error("The starting password and the hashed password don't match")
		return
	}
}
