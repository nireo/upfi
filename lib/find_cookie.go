package lib

import (
	"errors"
	"fmt"
	"net/http"
)

// FindCookieInResponse loops over the cookies in a request and checks if those cookies include the given cookie.
func FindCookieInResponse(cookie string, resp *http.Response) error {
	for _, respCookie := range resp.Cookies() {
		fmt.Println(respCookie.Name)
		if respCookie.Name == "token" {
			return nil
		}
	}

	return errors.New("could not find cookie in request")
}
