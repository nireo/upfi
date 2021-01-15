package templateapi

import (
	"net/http"
	"testing"

	"github.com/valyala/fasthttp"
)

// TestHomeRoute tests if going to the home page returns text/html and returns a successful status code
func TestHomeRoute(t *testing.T) {
	r, err := http.NewRequest("GET", "http://test/", nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := ServeRouter(CreateRouter().Handler, r)
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
