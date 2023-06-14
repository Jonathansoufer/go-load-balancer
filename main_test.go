package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimpleServer_Serve(t *testing.T) {
    server := newSimpleServer("http://localhost:8080")
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.Serve)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    expected := "Hello, World!"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestSimpleServer_IsAlive(t *testing.T) {
    server := newSimpleServer("http://localhost:8080")
    if !server.IsAlive() {
        t.Errorf("server should be alive")
    }
}