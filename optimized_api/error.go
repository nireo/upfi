package optimized_api

import "github.com/valyala/fasthttp"

// Define the fields that are on the error site.
type ErrorPageContent struct {
	StatusCode  int
	MainMessage string
	Description string
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
)
