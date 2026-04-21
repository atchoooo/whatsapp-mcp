package main

import (
	"net/http"
	"os"
)

// requireAPIKey rejects requests missing a valid X-API-Key header.
// If API_KEY env var is unset or empty, all requests are rejected (fail-secure).
func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	apiKey := os.Getenv("API_KEY")
	return func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" || r.Header.Get("X-API-Key") != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
