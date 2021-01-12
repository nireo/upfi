package jsonapi

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

// RunJSONAPI is just a simple functions that sets up the json router to run on a given port.
func RunJSONAPI(port string) {
	if err := fasthttp.ListenAndServe(fmt.Sprintf("localhost:%s", port), CreateJSONRouter().Handler); err != nil {
		log.Fatal(err)
	}
}
