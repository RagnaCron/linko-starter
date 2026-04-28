package main

import (
	"crypto/rand"
	"net/http"
)

func requestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = rand.Text()
			}
			w.Header().Set("X-Request-ID", reqID)

			next.ServeHTTP(w, r)
		})
	}
}
