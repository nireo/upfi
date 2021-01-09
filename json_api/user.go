package json_api

import (
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// WhoAmI handlers checks for a token and if a token was found return information about the user that has the token.
// This is used to set the user state in the front-end and have it up-to-date.
func WhoAmI(ctx *fasthttp.RequestCtx) {
	user, err := models.FindOneUser(&models.User{Username: string(ctx.Request.Header.Peek("username"))})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, user)
}
