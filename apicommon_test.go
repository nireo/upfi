package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nireo/upfi/templateapi"
	"github.com/valyala/fasthttp"
)

// TestAuthFilePagesReturnUnAuthorized checks if going to all the protected pages returns a unauthorized status code.
func TestAuthFilePagesReturnUnAuthorized(t *testing.T) {
	input := []struct {
		method string
		url    string
	}{
		{http.MethodGet, "upload"},
		{http.MethodGet, "file/awdawd"},
		{http.MethodGet, "settings"},
		{http.MethodGet, "files"},
		{http.MethodGet, "upload"},
	}

	for _, tt := range input {
		r, err := http.NewRequest(tt.method, fmt.Sprintf("http://test/%s", tt.url), nil)
		if err != nil {
			t.Error(err)
			return
		}

		res, err := templateapi.ServeRouter(templateapi.CreateRouter().Handler, r)
		if err != nil {
			t.Error(err)
			return
		}

		if res.StatusCode != fasthttp.StatusUnauthorized {
			t.Errorf("Wrong status code, wanted 401 got: %d", res.StatusCode)
		}
	}
}

// TestUploadPageLoadsWithToken tests if we append a auth token to the request and the handler returns text/html and a
// successful status code.
func TestPagesLoadWithToken(t *testing.T) {
	token, err := templateapi.NewTestUser("username", "password")
	if err != nil {
		t.Error(err)
		return
	}

	pages := []string{
		"upload",
		"files",
		"settings",
	}

	for _, page := range pages {
		r, err := http.NewRequest(http.MethodGet, "http://test/"+page, nil)
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

		if res.Header.Get("Content-Type") != "text/html" {
			t.Errorf("Wrong content type, wanted 'text/html' got: '%s'", res.Header.Get("Content-Type"))
			return
		}
	}
}
