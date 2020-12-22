package optimized_api

import (
	"html/template"

	"github.com/valyala/fasthttp"
)

func ErrorPageHandler(ctx *fasthttp.RequestCtx, errorType ErrorPageContent) {
	ctx.Response.SetStatusCode(errorType.StatusCode)
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/error_page.html"))
	// Execute the template giving it the errorType, so that it can properly display the
	// error message to the user
	if err := tmpl.Execute(ctx, errorType); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError)
		return
	}
}
