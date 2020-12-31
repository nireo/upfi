package main

import (
	"net/http"
	"testing"

	"github.com/nireo/upfi/optimized_api"
	"github.com/valyala/fasthttp"
)

func TestHomeRoute(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/", nil)
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
}
