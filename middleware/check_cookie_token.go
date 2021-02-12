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

// MoveIfAuthenticated is a middleware which is used when going to pages like login or register, since
// the user shouldn't go on the those pages if already authenticated.
func MoveIfAuthenticated(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		cookie := ctx.Request.Header.Cookie("token")
		if string(cookie) == "" {
			h(ctx)
			return
		}

		username, err := lib.ValidateToken(string(cookie))
		if err != nil {
			// since the token is not valid, the user isn't authenticated so we move the user
			// to the given handler.
			h(ctx)
			return
		}

		ctx.Request.Header.Add("username", username)
		ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
	}
}
