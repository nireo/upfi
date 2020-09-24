package optimized_api

import (
	"html/template"

	"github.com/valyala/fasthttp"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/not_found_template.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func ForbiddenHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/forbidden_template.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func InternalServerErrorHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/internal.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}
