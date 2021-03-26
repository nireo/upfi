package middleware

import (
	"github.com/nireo/upfi/lib"
	"github.com/valyala/fasthttp"
)

// CheckAuthentication looks for a cookie, given by the /register or /login routes. And finds the username
// in that jwt token.
func CheckAuthentication(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Take the cookie named token from the request headers.
		cookie := ctx.Request.Header.Cookie("token")

		// Use a function from the utils that verifies the integrity of a token and returns the
		// username in that token.
		username, err := lib.ValidateToken(string(cookie))
		if err == nil {
			// If there was no error, the token is valid and we can move on to the authenticated http handler.
			ctx.Request.Header.Add("username", username)
			h(ctx)
			return
		}

		// If the token validation failed we return a unauthorized status.
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	}
}
