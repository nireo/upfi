package templateapi

import (
	"html/template"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/templates"
	"github.com/valyala/fasthttp"
)

// ErrorPageHandler takes in a request context and a error type which is then used in a template to
// dynamically display an error page.
func ErrorPageHandler(ctx *fasthttp.RequestCtx, errorType lib.ErrorPageContent) {
	// Set the proper headers.
	ctx.Response.SetStatusCode(errorType.StatusCode)
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles(
		lib.AddRootToPath("templates/error_page.html")))
	// Execute the template giving it the errorType, so that it can properly display the
	// error message to the user
	if err := tmpl.Execute(ctx, errorType); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError)
		return
	}
}

// ServeHomePage just sends the user the home.html file which contains some information about the
// service.
func ServeHomePage(ctx *fasthttp.RequestCtx) {
	// Set the proper headers and then send the file.
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "text/html")
	// ctx.SendFile(lib.AddRootToPath("static/home.html"))
	templates.Home(ctx)
}

// RedirectToAuthorized is a handler that moves the user to an authorized page if logged in.
// For example: user goes to login page even though the user has an authorized token, so we move
// the user to the files page.
func RedirectToAuthorized(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}
