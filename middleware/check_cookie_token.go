package middleware

import (
	"github.com/nireo/upfi/lib"
	"github.com/valyala/fasthttp"
)

func CheckAuthentication(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		cookie := ctx.Request.Header.Cookie("auth")

		// validate token
		username, err := lib.ValidateToken(string(cookie))
		if err == nil {
			// token is validated successfully
			ctx.Request.Header.Add("username", username)
			h(ctx)
			return
		}

		// token validation failed
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	})
}
