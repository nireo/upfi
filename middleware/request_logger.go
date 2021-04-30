package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/valyala/fasthttp"
)

var (
	// Define an output such that the log package doesn't display the date (we don't want the date in every logger)
	output  = log.New(os.Stdout, "", 0)
	enabled = true
)

// SetHTTPLogging sets the variable in this file to a given condition
func SetHTTPLogging(condition bool) {
	enabled = condition
}

func getHTTPVersion(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// TinyLogger prints a few pieces of information to stdout when a request happens.
// FORMAT: [method] [url] [status]  [response-time]
func TinyLogger(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	if enabled {
		return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
			begin := time.Now()
			req(ctx)
			end := time.Now()
			output.Printf("%s %s %v %v", ctx.Method(), ctx.RequestURI(),
				ctx.Response.Header.StatusCode(), end.Sub(begin))
		})
	}

	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		req(ctx)
	})
}

// FullLogger prints out most of the provided information to stdout
// FORMAT [time] [remote-addr] [http-version] [method] [url] [status] [response-time] [user-agent]
func FullLogger(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	if enabled {
		return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
			begin := time.Now()
			req(ctx)
			end := time.Now()
			output.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
				end.Format("2006/01/02 - 15:04:05"),
				ctx.RemoteAddr(),
				getHTTPVersion(ctx),
				ctx.Method(),
				ctx.RequestURI(),
				ctx.Response.Header.StatusCode(),
				end.Sub(begin),
				ctx.UserAgent(),
			)
		})
	}

	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		req(ctx)
	})
}

// TinyLogger prints a few pieces of information to stdout when a request happens.
// FORMAT: [method] [url] [status]  [response-time]
func TinyHTTPLogger(req httprouter.Handle) httprouter.Handle {
	if enabled {
		return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			begin := time.Now()
			req(w, r, httprouter.Params{})
			end := time.Now()
			output.Printf("%s %s %v", r.Method, r.URL.RequestURI(), end.Sub(begin))
		})
	}

	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		req(w, r, httprouter.Params{})
	})
}
