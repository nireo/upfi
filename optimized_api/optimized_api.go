package optimized_api

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func SetupOptimizedApi() {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/files":
			// file handler
		default:
			fmt.Print("File not found")
		}
	}

	if err := fasthttp.ListenAndServe("localhost:8080", requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}
