package crypt

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCrypt(t *testing.T) {
	buf, err := ioutil.ReadFile("./test_file.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if err := EncryptToDst("./encrypted.txt", buf, "test"); err != nil {
		t.Error(err)
		return
	}

	// compare the bytes in the files
	encryptedBuf, err := ioutil.ReadFile("./encrypted.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Equal(encryptedBuf, buf) {
		t.Error("Bytes are the same even though encrypted")
		return
	}

	// decrypt the file
	if err := DecryptToDst("./decrypted.txt", "./encrypted.txt", "test"); err != nil {
		t.Error(err)
		return
	}

	decryptedBuf, err := ioutil.ReadFile("./decrypted.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(decryptedBuf, buf) {
		t.Error("Bytes are not the same even though data decrypted")
		return
	}

	// Remove all the extra files in the end
	if err := os.Remove("./decrypted.txt"); err != nil {
		t.Error(err)
		return
	}

	if err := os.Remove("./encrypted.txt"); err != nil {
		t.Error(err)
		return
	}
}
