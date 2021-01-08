package lib

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func WriteResponseJSON(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.SetStatusCode(code)
	ctx.Response.Header.Set("Content-Type", "application/json")
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}