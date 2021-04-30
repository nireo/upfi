package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/valyala/fasthttp"
)

var (
	corsAllowHeader      = "authorization,content-type"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

// CORS middleware allows us to make requests to the seperate javascript client
func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeader)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func HTTPCORS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		w.Header().Set("Access-Control-Allow-Headers", corsAllowHeader)
		w.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
	}
}
