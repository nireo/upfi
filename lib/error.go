package lib

import "github.com/valyala/fasthttp"

// Define the fields that are on the error site.
type ErrorPageContent struct {
	StatusCode  int    `json:"status"`
	MainMessage string `json:"message"`
	Description string `json:"description"`
}

var (
	InternalServerErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusInternalServerError,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
		Description: "There was a server error while processing your request, this is due to the server.",
	}

	BadRequestErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusBadRequest,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		Description: "There was an error without request input.",
	}

	NotFoundErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusNotFound,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		Description: "The thing you're searching for has not been found",
	}

	ForbiddenErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusForbidden,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusForbidden),
		Description: "You're not allowed to view the content on this page.",
	}

	ConflictErrorPage = ErrorPageContent{
		StatusCode:  fasthttp.StatusConflict,
		MainMessage: fasthttp.StatusMessage(fasthttp.StatusConflict),

		// This isn't really a 'common' description, but in the case of this application
		// the conflict error only happens during registration.
		Description: "The username is already taken.",
	}
)

func CreateSimpleErrorContent(statusCode int) *ErrorPageContent {
	return &ErrorPageContent{
		StatusCode:  statusCode,
		MainMessage: fasthttp.StatusMessage(statusCode),
		Description: "No description provided.",
	}
}
