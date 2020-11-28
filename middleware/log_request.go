package middleware

import (
	"log"
	"net/http"
	"time"
)

// LogRequest takes a http.HandlerFunc as a parameter and returns the handler functions wrapped
// around a logger, which logs the execution time and the path of a request.
func LogRequest() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()
			f(w, r)
		}
	}
}
