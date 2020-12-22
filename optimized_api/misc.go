package optimized_api

import (
	"html/template"

	"github.com/valyala/fasthttp"
)

// NotFoundHandler is a http handler which takes the request context as an input, and then uses
// that context to display a not found page. This is probably worse performance wise, or insignificant, but
// it makes the code cleaner.
func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	// Set the right status code
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
	// Set the right content type, so the html template can be properly viewed.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Execute a html file as a template and write that template, using the request context.
	tmpl := template.Must(template.ParseFiles("./templates/not_found_template.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// ForbiddenHandler is a http handler which takes a request context as input, and uses the context, to
// return a forbidden page using a html template.
func ForbiddenHandler(ctx *fasthttp.RequestCtx) {
	// Set the right status code
	ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
	// Set the right content type, so the html template can be properly viewed.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Execute a html file as a template and write that template, using the request context.
	tmpl := template.Must(template.ParseFiles("./templates/forbidden_template.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// InternalServerErrorHandler is a http handler which takes a request context as input, and uses the context, to
// return a Internal Server Error page using a html template.
func InternalServerErrorHandler(ctx *fasthttp.RequestCtx) {
	// Set the right status code
	ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
	// Set the right content type, so the html template can be properly viewed.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Execute a html file as a template and write that template, using the request context.
	tmpl := template.Must(template.ParseFiles("./templates/internal.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// InternalServerErrorHandler is a http handler which takes a request context as input, and uses the context, to
// return a Internal Server Error page using a html template.
func BadRequestHandler(ctx *fasthttp.RequestCtx) {
	// Set the right status code
	ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
	// Set the right content type, so the html template can be properly viewed.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Execute a html file as a template and write that template, using the request context.
	tmpl := template.Must(template.ParseFiles("./templates/bad_request.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}
