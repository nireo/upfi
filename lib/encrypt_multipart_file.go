package lib

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/nireo/upfi/crypt"
)

// EncryptMultipartFile encrypts a file into a location and then returns the file's MIME address,
// which is needed to properly let the user download content.
func EncryptMultipartFile(file multipart.File, path, key string) (string, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	if err := crypt.EncryptToDst(path, buf.Bytes(), key); err != nil {
		return "", err
	}

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}

	return http.DetectContentType(fileHeader), nil
}
