package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/nireo/booru/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/optimized_api"
	"github.com/valyala/fasthttp"
)

func TestLoginRoute(t *testing.T) {
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

func TestRegisterRoute(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/register", nil)
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

// TestRegister first tests if account creation works through http and then tests
// if removing user using different helper functions works.
func TestRegister(t *testing.T) {
	db := lib.GetDatabase()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("master", "secret")
	_ = writer.WriteField("username", "testaccount12345")
	_ = writer.WriteField("password", "reallysecretpassword")

	if err := writer.Close(); err != nil {
		t.Error(err)
	}

	r, err := http.NewRequest(fasthttp.MethodPost, "http://test/register", body)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Errorf("Wrong status code, wanted 200 got: %d", res.StatusCode)
	}

	if res.StatusCode != fasthttp.StatusOK {
		t.Errorf("Wrong status code, wanted 200 got: %d", res.StatusCode)
	}

	// check that the user has been created
	var user models.User
	if err := db.Where(&models.User{Username: "testaccount12345"}).First(&user).Error; err != nil {
		t.Error("User was not created, err: ", err)
	}

	// check that a folder has been created
	if _, err := os.Stat("./files/" + user.UUID); os.IsNotExist(err) {
		t.Error("A file folder wasn't created for the user", err)
	}

	// after all this remove the user
	if err := user.Delete(); err != nil {
		t.Error("Could not remove user, err: ", err)
	}
}
