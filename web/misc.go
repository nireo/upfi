package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/templates"
)

// ErrorPageHandler takes in a request context and a error type which is then used in a template to
// dynamically display an error page.
func ErrorPageHandler(w http.ResponseWriter, r *http.Request,
	errorType lib.ErrorPageContent) {
	// Set the proper headers.
	w.WriteHeader(errorType.StatusCode)
	w.Header().Set("Content-Type", "text/html")

	params := templates.ErrorParams{
		Title:         errorType.MainMessage,
		Authenticated: lib.IsAuth(r),
		Error:         errorType,
	}

	templates.ErrorPage(w, params)
}

// ServeHomePage just sends the user the home.html file which contains some information about the
// service.
func ServeHomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	templates.Home(w, templates.HomeParams{
		Title:         "home",
		Authenticated: lib.IsAuth(r),
	})
}
