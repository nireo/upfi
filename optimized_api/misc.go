package optimized_api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	// prompt the user
	fmt.Fprintf(ctx, "Cannot: '%s' route: '%s'", ctx.Method(), ctx.RequestURI())
	ctx.SetContentType("text/plain; charset=utf-8")

	// set not found status
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
}
