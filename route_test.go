package muxify

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Fatal("Expected NewRouter() to return a non-nil router")
	}
}

func TestRouteRegistrationAndHandling(t *testing.T) {
	router := NewRouter()
	expectedResponse := "Hello, World!"

	router.Handle("GET", "/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse)
	}))

	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Body.String() != expectedResponse {
		t.Errorf("Expected response to be '%s', but got '%s'", expectedResponse, recorder.Body.String())
	}
}

func TestMiddlewareExecution(t *testing.T) {
	router := NewRouter()
	expectedResponse := "Hello, Middleware!"

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expectedResponse)
			next.ServeHTTP(w, r)
		})
	})

	router.Handle("GET", "/middleware", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req, err := http.NewRequest("GET", "/middleware", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Body.String() != expectedResponse {
		t.Errorf("Expected response to be '%s', but got '%s'", expectedResponse, recorder.Body.String())
	}
}
