package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// EncryptToDst takes in a destination path, file data and an encryption key. The given data is encrypted
// using the encryption key and written into the dst location.
func EncryptToDst(dst string, data []byte, key string) error {
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	encryptedData, err := encrypt(data, key)
	if err != nil {
		return err
	}
	file.Write(encryptedData)

	return nil
}
