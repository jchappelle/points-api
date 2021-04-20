package web

import (
	"log"
	"net/http"
)

// loggingMiddleware logs the method and url of each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request: %s %s", r.Method, r.URL.EscapedPath())
			next.ServeHTTP(w, r)
		},
	)
}