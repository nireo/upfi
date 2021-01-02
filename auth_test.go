package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/optimized_api"
	"github.com/valyala/fasthttp"
)

func TestNewTestUser(t *testing.T) {
	username, password := "username", "password"
	token, err := optimized_api.NewTestUser(username, password)
	if err != nil {
		t.Error(err)
		return
	}

	if token == "" {
		t.Error("The jwt token is empty")
		return
	}

	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		t.Error("Could not find user, err: ", err)
	}

	if _, err := os.Stat("./files/" + user.UUID); os.IsNotExist(err) {
		t.Error("A file folder wasn't created for the user", err)
	}

	if err := user.Delete(); err != nil {
		t.Error("Could not remove user, err: ", err)
	}
}

func TestLoginRouteGet(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/login", nil)
	if err != nil {
		t.Error(err)
	}

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != fasthttp.StatusOK {
		t.Errorf("Wrong status code, wanted 200 got: %d", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "text/html" {
		t.Errorf("Wrong Content-Type, wanted 'text/html' got: %s", res.Header.Get("Content-Type"))
	}
}

func TestRegisterRouteGet(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/register", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != fasthttp.StatusOK {
		t.Errorf("Wrong status code, wanted 200 got: %d", res.StatusCode)
		return
	}

	if res.Header.Get("Content-Type") != "text/html" {
		t.Errorf("Wrong Content-Type, wanted 'text/html' got: %s", res.Header.Get("Content-Type"))
		return
	}
}

// TestRegister first tests if account creation works through http and then tests
// if removing user using different helper functions works.
func TestRegister(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("master", "secret")
	_ = writer.WriteField("username", "user")
	_ = writer.WriteField("password", "reallysecretpassword")

	if err := writer.Close(); err != nil {
		t.Error(err)
		return
	}

	r, err := http.NewRequest(fasthttp.MethodPost, "http://test/register", body)
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())

	_, err = optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	// check that the user has been created
	user, err := models.FindOneUser(&models.User{Username: "user"})
	if err != nil {
		t.Error("Could not find user, err: ", err)
		return
	}

	// check that a folder has been created
	if _, err := os.Stat("./files/" + user.UUID); os.IsNotExist(err) {
		t.Error("A file folder wasn't created for the user", err)
		return
	}

	// after all this remove the user
	if err := user.Delete(); err != nil {
		t.Error("Could not remove user, err: ", err)
		return
	}
}

func TestRegisterInvalidInput(t *testing.T) {
	// Without sending any data
	r, err := http.NewRequest(fasthttp.MethodPost, "http://test/register", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf("Wrong status code, wanted 400 got: %d", res.StatusCode)
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("master", "secret")
	_ = writer.WriteField("username", "user")
	_ = writer.WriteField("password", "reallysecretpassword")

	if err := writer.Close(); err != nil {
		t.Error(err)
		return
	}

	r2, err := http.NewRequest(fasthttp.MethodPost, "http://test/register", body)
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())

	res2, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r2)
	if err != nil {
		t.Error(err)
		return
	}

	if res2.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf("Wrong status code, wanted 401 got: %d", res.StatusCode)
		return
	}
}

// TestLoginInvalidInput tests for problems with different types of edge cases for
// the login route.
func TestLoginInvalidInput(t *testing.T) {
	// Without sending any data
	r, err := http.NewRequest(fasthttp.MethodPost, "http://test/login", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf("Wrong status code, wanted 400 got: %d", res.StatusCode)
		return
	}
}
