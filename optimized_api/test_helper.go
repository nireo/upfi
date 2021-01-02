package optimized_api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func ServeRouter(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

// NewTestUser takes in a username and password, such that testing with authenticated routes is easier.
// does pretty much everything like the normal register route except use http. Returns a token and an error.
func NewTestUser(username, password string) (string, error) {
	if len(username) < 3 || len(password) < 8 {
		return "", errors.New("password and/or username too short")
	}

	if len(username) > 20 || len(password) > 32 {
		return "", errors.New("password and/or username too long")
	}

	// no need to hash since it takes a lot of time and the user will be deleted after the tests.
	newUser := models.User{
		Username:             username,
		Password:             password,
		FileEncryptionMaster: "secret",
		UUID:                 lib.GenerateUUID(),
	}

	if err := os.Mkdir("./files/"+newUser.UUID, os.ModePerm); err != nil {
		return "", errors.New("could not create an user file directory")
	}

	db := lib.GetDatabase()
	db.Create(&newUser)

	token, err := lib.CreateToken(newUser.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateMultipartForm(fields map[string]string) (*bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
