package crypt

import (
	"bytes"
	"github.com/nireo/upfi/lib"
	"io/ioutil"
	"os"
	"testing"
)

func TestCrypt(t *testing.T) {
	buf, err := ioutil.ReadFile(lib.AddRootToPath("crypt/test_file.txt"))
	if err != nil {
		t.Error(err)
		return
	}

	decPath := lib.AddRootToPath("crypt/encrypted.txt")
	encPath := lib.AddRootToPath("crypt/decrypted.txt")

	if err := EncryptToDst(encPath, buf, "test"); err != nil {
		t.Error(err)
		return
	}


	// compare the bytes in the files
	encryptedBuf, err := ioutil.ReadFile(encPath)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Equal(encryptedBuf, buf) {
		t.Error("Bytes are the same even though encrypted")
		return
	}

	// decrypt the file
	if err := DecryptToDst(decPath, encPath, "test"); err != nil {
		t.Error(err)
		return
	}

	decryptedBuf, err := ioutil.ReadFile(decPath)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(decryptedBuf, buf) {
		t.Error("Bytes are not the same even though data decrypted")
		return
	}

	// Remove all the extra files in the end
	if err := os.Remove(decPath); err != nil {
		t.Error(err)
		return
	}

	if err := os.Remove(encPath); err != nil {
		t.Error(err)
		return
	}
}
