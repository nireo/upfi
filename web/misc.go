package web

import (
	"net/http"

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
