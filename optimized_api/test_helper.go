package optimized_api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func SetupServerWithHandler(port, method string, handler fasthttp.RequestHandler) {
	router := fasthttprouter.New()

	// Since when we use this helper the we're testing a single route. Using '/' also makes the test code
	// more clear.
	router.Handle(method, "/", handler)
	if err := fasthttp.ListenAndServe(fmt.Sprintf("localhost:%s", port), router.Handler); err != nil {
		log.Fatal("Error in ListenAndServe")
	}
}

func ServeRouter(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}
