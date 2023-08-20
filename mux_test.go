package muxify

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	mux := NewMux()
	if mux == nil {
		t.Fatal("Expected NewMux() to return a non-nil Mux")
	}
}

func TestMuxRouteRegistrationAndHandling(t *testing.T) {
	mux := NewMux()
	expectedResponse := "Hello, World!"

	mux.Handle("GET", "/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse)
	}))

	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	if recorder.Body.String() != expectedResponse {
		t.Errorf("Expected response to be '%s', but got '%s'", expectedResponse, recorder.Body.String())
	}
}

func TestMuxMiddlewareExecution(t *testing.T) {
	mux := NewMux()
	expectedResponse := "Hello, Middleware!"

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expectedResponse)
			next.ServeHTTP(w, r)
		})
	})

	mux.Handle("GET", "/middleware", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req, err := http.NewRequest("GET", "/middleware", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	if recorder.Body.String() != expectedResponse {
		t.Errorf("Expected response to be '%s', but got '%s'", expectedResponse, recorder.Body.String())
	}
}

func TestMuxNotFoundHandling(t *testing.T) {
	mux := NewMux()
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code to be '%d', but got '%d'", http.StatusNotFound, recorder.Code)
	}
}
