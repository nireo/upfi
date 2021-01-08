package json_api

import (
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// WhoAmI handlers checks for a token and if a token was found return information about the user that has the token
func WhoAmI(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, user)
}