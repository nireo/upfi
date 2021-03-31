package lib

import "github.com/valyala/fasthttp"

// ErrorPageContent defines the fields that are on the error site.
type ErrorPageContent struct {
	StatusCode  int    `json:"status"`
	MainMessage string `json:"message"`
	Description string `json:"description"`
}

var (
	// InternalServerErrorPage is used when some internal code execution has failed.
	InternalServerErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusInternalServerError,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
		Description: "There was a server error while processing your request, this is due to the server.",
	}

	// BadRequestErrorPage is used when the user's request had bad properties.
	BadRequestErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusBadRequest,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		Description: "There was an error without request input.",
	}

	// NotFoundErrorPage is used when something the user is searching for is not found.
	NotFoundErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusNotFound,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		Description: "The thing you're searching for has not been found",
	}

	// ForbiddenErrorPage is used when the user is trying to do something without permission.
	ForbiddenErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusForbidden,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusForbidden),
		Description: "You're not allowed to view the content on this page.",
	}

	// ConflictErrorPage is used when the user tries to create information into the database that
	// already exists.
	ConflictErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusConflict,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusConflict),

		// This isn't really a 'common' description, but in the case of this application
		// the conflict error only happens during registration.
		Description: "The username is already taken.",
	}
)

// CreateSimpleErrorContent makes a simpler and more common version of the above constants. This is
// used for status codes that are not as common as the ones above.
func CreateSimpleErrorContent(statusCode int) *ErrorPageContent {
	return &ErrorPageContent{
		StatusCode:  statusCode,
		MainMessage: fasthttp.StatusMessage(statusCode),
		Description: "No description provided.",
	}
}

// CreateDetailedErrorContent makes a more detailed  error page with fully custom content
func CreateDetailedErrorContent(err error, title string, code int) *ErrorPageContent {
	return &ErrorPageContent{
		StatusCode:  code,
		MainMessage: title,
		Description: err.Error(),
	}
}
