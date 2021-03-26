package lib

import "github.com/valyala/fasthttp"

// IsAuthenticated checks if there is a token the request context
func IsAuthenticated(ctx *fasthttp.RequestCtx) bool {
	cookie := ctx.Request.Header.Cookie("token")
	_, err := ValidateToken(string(cookie))
	if err == nil {
		return true
	}

	return false
}
