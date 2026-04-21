package main

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// requireAPIKey rejects requests missing a valid X-API-Key header.
// If API_KEY env var is unset or empty, all requests are rejected (fail-secure).
func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	apiKey := os.Getenv("API_KEY")
	return func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get("X-API-Key")), []byte(apiKey)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
