package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))

	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return key, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return key, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return plaintext, err
	}

	return plaintext, nil
}

func DecryptToDst(dst, src, key string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	decryptedData, err := decrypt(data, key)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dst, decryptedData, 0666); err != nil {
		return err
	}

	return nil
}
