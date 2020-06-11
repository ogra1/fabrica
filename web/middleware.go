package web

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Logger Handle logging for the web service
func Logger(start time.Time, r *http.Request) {
	// Reduce noise in logs
	if r.Method == "GET" && r.RequestURI == "/v1/system" {
		return
	}
	if r.Method == "GET" && strings.HasPrefix(r.RequestURI, "/v1/builds/") {
		return
	}

	log.Printf(
		"%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		time.Since(start),
	)
}

// Middleware to pre-process web service requests
func Middleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request
		Logger(start, r)

		inner.ServeHTTP(w, r)
	})
}
