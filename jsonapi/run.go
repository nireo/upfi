package jsonapi

import (
	"fmt"
	"log"

	"github.com/nireo/upfi/middleware"
	"github.com/valyala/fasthttp"
)

// RunJSONAPI is just a simple functions that sets up the json router to run on a given port.
func RunJSONAPI(port string) {
	if err := fasthttp.ListenAndServe(fmt.Sprintf("localhost:%s", port),
		middleware.CORS(CreateJSONRouter().Handler)); err != nil {
		log.Fatal(err)
	}
}
