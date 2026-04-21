package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAPIKey_NoHeader(t *testing.T) {
	t.Setenv("API_KEY", "test-secret")
	handler := requireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/api/send", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIKey_WrongKey(t *testing.T) {
	t.Setenv("API_KEY", "test-secret")
	handler := requireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/api/send", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIKey_CorrectKey(t *testing.T) {
	t.Setenv("API_KEY", "test-secret")
	handler := requireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/api/send", nil)
	req.Header.Set("X-API-Key", "test-secret")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRequireAPIKey_EmptyEnvAlwaysRejects(t *testing.T) {
	t.Setenv("API_KEY", "")
	handler := requireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/api/send", nil)
	req.Header.Set("X-API-Key", "")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("empty API_KEY should reject all requests, got %d", w.Code)
	}
}
