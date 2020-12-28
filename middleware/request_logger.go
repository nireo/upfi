package middleware

import (
	"log"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	// Define an output such that the log package doesn't display the date (we don't want the date in every logger)
	output = log.New(os.Stdout, "", 0)
)

func getHttpVersion(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// TinyLogger prints a few pieces of information to stdout when a request happens.
// FORMAT: [method] [url] [status]  [response-time]
func TinyLogger(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		log.Printf("%s %s %v %v", ctx.Method(), ctx.RequestURI(),
			ctx.Response.Header.StatusCode(), end.Sub(begin))
	})
}

// FullLogger prints out most of the provided information to stdout
// FORMAT [time] [remote-addr] [http-version] [method] [url] [status] [response-time] [user-agent]
func FullLogger(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		log.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			getHttpVersion(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}
