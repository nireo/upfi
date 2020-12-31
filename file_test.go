package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nireo/upfi/optimized_api"
)

// TestFileUploadPageForbidden tests that when we try to access a site without authorization
// it returns a status code.
func TestFileUploadPageUnAuthorized(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://test/upload", nil)
	if err != nil {
		t.Error(err)
	}

	res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Wrong status code, wanted 401 got: %d", res.StatusCode)
	}
}

func AuthFilePagesReturnUnAuthorized(t *testing.T) {
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
		}

		res, err := optimized_api.ServeRouter(optimized_api.CreateRouter().Handler, r)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != http.StatusUnauthorized {
			t.Errorf("Wrong status code, wanted 401 got: %d", res.StatusCode)
		}
	}
}
