package crypt

import (
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

// DeriveKey takes in a password and a salt, then returns the derived key and the password.
func DeriveKey(password string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}

	return pbkdf2.Key([]byte(password), salt, 1000, 32, sha256.New), salt
}
