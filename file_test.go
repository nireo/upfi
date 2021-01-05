package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nireo/upfi/optimized_api"
	"github.com/valyala/fasthttp"
)

// TestFileUploadPageForbidden tests that when we try to access a site without authorization
// it returns a status code.
func TestFileUploadPageUnAuthorized(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://test/upload", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != fasthttp.StatusUnauthorized {
		t.Errorf("Wrong status code, wanted 401 got: %d", res.StatusCode)
		return
	}
}

// TestAuthFilePagesReturnUnAuthorized checks if going to all the protected pages returns a unauthorized status code.
func TestAuthFilePagesReturnUnAuthorized(t *testing.T) {
	toTest := []string{
		"upload",
		"file/testsetsetse",
		"files",
		"settings",
	}

	for _, page := range toTest {
		r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://test/%s", page), nil)
		if err != nil {
			t.Error(err)
			return
		}

		res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
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
func TestUploadPageLoadsWithToken(t *testing.T) {
	token, err := optimized_api.NewTestUser("username", "password")
	if err != nil {
		t.Error(err)
		return
	}

	r, err := http.NewRequest(http.MethodGet, "http://test/upload", nil)
	if err != nil {
		t.Error(err)
		return
	}

	r.AddCookie(&http.Cookie{Name: "token", Value: token})
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
		t.Errorf("Wrong content type, wanted 'text/html' got: '%s'", res.Header.Get("Content-Type"))
		return
	}
}
