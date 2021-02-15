package main

import (
	"bytes"
	"fmt"
	"github.com/nireo/upfi/lib"
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

	if _, err := os.Stat(lib.AddRootToPath("files/") + user.UUID); os.IsNotExist(err) {
		t.Error("A file folder wasn't created for the user", err)
	}

	if err := user.Delete(); err != nil {
		t.Error("Could not remove user, err: ", err)
	}
}

// TestGetPages checks if there is html returned for the 'login' and 'register' pages.
func TestGetPages(t *testing.T) {
	tests := []string{
		"login",
		"register",
	}

	for _, testRoute := range tests {
		r, err := http.NewRequest("GET", "http://test/"+testRoute, nil)
		if err != nil {
			t.Error(err)
			continue
		}

		res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
		if err != nil {
			t.Errorf("error creating request to %s, err: %s", testRoute, err)
			continue
		}

		if res.StatusCode != fasthttp.StatusOK {
			t.Errorf("got the wrong status code for %s. want=200, got=%d", testRoute, res.StatusCode)
			continue
		}

		if res.Header.Get("Content-Type") != "text/html" {
			t.Errorf("Wrong Content-Type for %s, wanted 'text/html', got=%s",
				testRoute, res.Header.Get("Content-Type"))
		}
	}

}

// TestRegister first tests if account creation works through http and then tests
// if removing user using different helper functions works.
func TestRegisterComprehensive(t *testing.T) {
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

	if _, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r); err != nil {
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
	if _, err := os.Stat(lib.AddRootToPath("files/") + user.UUID); os.IsNotExist(err) {
		t.Error("A file folder wasn't created for the user", err)
		return
	}

	// after all this remove the user
	if err := user.Delete(); err != nil {
		t.Error("Could not remove user, err: ", err)
		return
	}
}

func TestRegisterRouteInputs(t *testing.T) {
	// define test cases like this so we don't need multiple functions to basically tests the same thing.
	testCases := []struct {
		//[0] username | [1] password | [2] master
		inputs         []string
		expectedStatus int
	}{
		{inputs: []string{"", "", ""}, expectedStatus: fasthttp.StatusBadRequest},
		{inputs: []string{"registertest1", "pas", "2short"}, expectedStatus: fasthttp.StatusBadRequest},
		{inputs: []string{"registertest2", "2short", "reallysecretpas"}, expectedStatus: fasthttp.StatusBadRequest},
		{inputs: []string{"2s", "secretpass", "secretpass"}, expectedStatus: fasthttp.StatusBadRequest},

		// These are valid request inputs, but the expectedStatus is 401, because the user is redirected to the
		// /files page, and since the request doesn't have the authorization cookie, the user cannot access the
		// page. But if the request shows this status it means the registeration process worked according to plan.
		{inputs: []string{"registertest3", "secretpass", "secretpass"}, expectedStatus: fasthttp.StatusUnauthorized},
	}

	for testNum, testCase := range testCases {
		// Create request body from test case inputs
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("username", testCase.inputs[0])
		_ = writer.WriteField("password", testCase.inputs[1])
		_ = writer.WriteField("master", testCase.inputs[2])

		if err := writer.Close(); err != nil {
			t.Error(err)
			return
		}

		// Create & send the request to the server.
		r, err := http.NewRequest(fasthttp.MethodPost, "http://test/register", body)
		if err != nil {
			t.Error(err)
			return
		}
		r.Header.Set("Content-Type", writer.FormDataContentType())

		res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
		if err != nil {
			t.Error(err)
			return
		}

		if res.StatusCode != testCase.expectedStatus {
			t.Errorf("wrong statuscode on test %d. want=%d, got=%d", testNum+1, testCase.expectedStatus, res.StatusCode)
			return
		}
	}
}

func TestLoginInputs(t *testing.T) {
	// create a new user so we can test if logging in is possible
	username, password := "realuser123", "password123"
	if _, err := templateapi.NewTestUser(username, password); err != nil {
		t.Error(err)
		return
	}

	testCases := []struct {
		//[0] username | [1] password
		inputs         []string
		expectedStatus int
	}{
		{inputs: []string{"", ""}, expectedStatus: fasthttp.StatusBadRequest},

		// This request is the only valid request, and the unauthorized status just means we were moved
		// to the /files page and we didn't have a token. But it still means that the login process worked.
		{inputs: []string{"realuser123", "password123"}, expectedStatus: fasthttp.StatusUnauthorized},
		{inputs: []string{"realuser123", "321password"}, expectedStatus: fasthttp.StatusForbidden},
		{inputs: []string{"notrealuser123", "secretpass"}, expectedStatus: fasthttp.StatusNotFound},
	}

	for testNum, testCase := range testCases {
		// Create request body from test case inputs
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("username", testCase.inputs[0])
		_ = writer.WriteField("password", testCase.inputs[1])

		if err := writer.Close(); err != nil {
			t.Error(err)
			return
		}

		// Create & send the request to the server.
		r, err := http.NewRequest(fasthttp.MethodPost, "http://test/login", body)
		if err != nil {
			t.Error(err)
			return
		}
		r.Header.Set("Content-Type", writer.FormDataContentType())

		res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
		if err != nil {
			t.Error(err)
			return
		}

		if res.StatusCode != testCase.expectedStatus {
			t.Errorf("wrong statuscode on test %d. want=%d, got=%d", testNum+1, testCase.expectedStatus, res.StatusCode)
			return
		}
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
