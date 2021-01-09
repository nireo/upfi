package json_api

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

func RunJSONApi(port string) {
	if err := fasthttp.ListenAndServe(fmt.Sprintf("localhost:%s", port), CreateJSONRouter().Handler); err != nil {
		log.Fatal(err)
	}
}

