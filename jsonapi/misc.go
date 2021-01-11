package jsonapi

import (
	"github.com/nireo/upfi/lib"
	"github.com/valyala/fasthttp"
)

// ServeErrorJSON takes in the request context and error type such that the user is returned a
// proper error message instead of the normal message. This handler is used by the front-end to
// properly display an error notification.
func ServeErrorJSON(ctx *fasthttp.RequestCtx, errorType lib.ErrorPageContent) {
	lib.WriteResponseJSON(ctx, errorType.StatusCode, errorType)
}
