package lib

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

// WriteResponseJSON takes in a request context, status code, and a golang interface. The functions sets
// the request's response with the given status code and content type. And then golang's json package
// encodes the data from the interface into json format.
func WriteResponseJSON(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.SetStatusCode(code)
	ctx.Response.Header.Set("Content-Type", "application/json")
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
