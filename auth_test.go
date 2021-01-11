package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/templateapi"
	"github.com/valyala/fasthttp"
)

// TestNewTestUser creates a user using the NewTestUser function which is used in tests to easily create a user.
func TestNewTestUser(t *testing.T) {
	username, password := "username", "password"
	token, err := templateapi.NewTestUser(username, password)
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

// TestLoginRouteGet tests if going to the login handler page returns a text/html content page with a successful status
// code.
func TestLoginRouteGet(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/login", nil)
	if err != nil {
		t.Error(err)
	}

	res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
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

// TestRegisterRouteGet tests if going to the register handler page returns a text/html content page with a successful
// status code.
func TestRegisterRouteGet(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/register", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
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

	_ = writer.WriteField("master", "secret123456")
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

	_, err = templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
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

// TestRegisterInvalidInput tests the register handler with different edge case input that should be handeled properly
// and the handler should return the fasthttp.StatusBadRequest which is 400.
func TestRegisterInvalidInput(t *testing.T) {
	// Without sending any data
	r, err := http.NewRequest(fasthttp.MethodPost, "http://test/register", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
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

	res2, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r2)
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

	res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf("Wrong status code, wanted 400 got: %d", res.StatusCode)
		return
	}
}

// TestLoginRoutePost tests a login into the login page. We add a token to the request ourselves, since for some reason
// the http.Response doesn't include a token. The fasthttp.StatusOK check if for testing if the user was successfully
// redirected to the files page.
func TestLoginRoutePost(t *testing.T) {
	// create a new user
	token, err := templateapi.NewTestUser("logintest", "password")
	if err != nil {
		t.Error(err)
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("username", "logintest")
	_ = writer.WriteField("password", "password")

	if err := writer.Close(); err != nil {
		t.Error(err)
		return
	}

	r, err := http.NewRequest(fasthttp.MethodPost, "http://test/login", body)
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.AddCookie(&http.Cookie{Name: "token", Value: token})

	res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	// We are directed to the /files page, so we check if redirecting us was successful.
	if res.StatusCode != fasthttp.StatusOK {
		t.Errorf("Wrong status code, wanted 200 got: %d", res.StatusCode)
		return
	}

	// check that the content type matches that of the files page.
	if res.Header.Get("Content-Type") != "text/html" {
		t.Errorf("Wrong content type, wanted 'text/html' got: %s", res.Header.Get("Content-Type"))
	}
}

// RedirectIfToken checks if the routes '/login' and '/register' redirect when the user enters them when owning a auth
// token. The routes should return the status code 308 which stands for MovedPermanently.
func TestRedirectIfToken(t *testing.T) {
	token, err := templateapi.NewTestUser("redirecttest", "password")
	if err != nil {
		t.Error(err)
		return
	}

	routesToCheck := []string{
		"login",
		"register",
	}

	for _, route := range routesToCheck {
		r, err := http.NewRequest(fasthttp.MethodPost, fmt.Sprintf("http://test/%s", route), nil)
		if err != nil {
			t.Error(err)
			return
		}
		r.AddCookie(&http.Cookie{Name: "token", Value: token})

		res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
		if err != nil {
			t.Error(err)
			return
		}

		if res.StatusCode != fasthttp.StatusOK {
			t.Errorf("Wrong status code, wanted 200 got: %d", res.StatusCode)
			return
		}

		if res.Request.URL.String() != "http://test/files" {
			t.Errorf("Wrong ending url, wanted 'http://test/files' got: %s", res.Request.URL.String())
			return
		}
	}
}