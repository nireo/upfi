package lib

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// IsAuthenticated checks if there is a token the request context
func IsAuthenticated(ctx *fasthttp.RequestCtx) bool {
	cookie := ctx.Request.Header.Cookie("token")
	_, err := ValidateToken(string(cookie))
	if err == nil {
		return true
	}

	return false
}

func IsAuth(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	if _, err := ValidateToken(string(cookie.Value)); err == nil {
		return true
	}

	return false
}
