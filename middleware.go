package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request method and URL
		log.Printf("Method: %s, URL: %s\n", r.Method, r.URL.String())

		// Log headers
		log.Println("Headers:")
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("  %s: %s\n", name, value)
			}
		}

		// Log query parameters
		log.Println("Query Parameters:")
		queryParams := r.URL.Query()
		for key, values := range queryParams {
			for _, value := range values {
				log.Printf("  %s: %s\n", key, value)
			}
		}

		// Log request body (if applicable)
		if r.Body != nil {
			// Read the body and replace it so the handler can still access it
			var buf bytes.Buffer
			tee := io.TeeReader(r.Body, &buf)
			body, err := io.ReadAll(tee)
			if err != nil {
				log.Printf("Failed to read body: %v\n", err)
			} else {
				log.Printf("Body: %s\n", string(body))
			}
			r.Body = io.NopCloser(&buf)
		}
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
