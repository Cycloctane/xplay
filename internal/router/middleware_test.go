package router

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("200"))
	w.WriteHeader(200)
})

func testBasicAuth(t *testing.T, username, password string, expectedStatus int) {
	logger := log.New(io.Discard, "", 0)
	handler := NewBasicAuth(testHandler, logger, "yes", "yes")
	req := httptest.NewRequest("GET", "/", nil)
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != expectedStatus {
		t.Errorf("%s:%s - Expected status %d, got %d", username, password, expectedStatus, rec.Code)
	}
}

func TestBasicAuth(t *testing.T) {
	t.Run("valid credentials", func(t *testing.T) {
		testBasicAuth(t, "yes", "yes", 200)
	})
	t.Run("invalid credentials", func(t *testing.T) {
		testBasicAuth(t, "no", "no", 401)
		testBasicAuth(t, "yes", "", 401)
		testBasicAuth(t, "", "", 401)
		testBasicAuth(t, "", "yes", 401)
	})
}

func testHostHeader(t *testing.T, host string, expectedStatus int) {
	handler := NewHostHeaderValidator(testHandler, "yes.local,127.0.0.1,example.com")
	req := httptest.NewRequest("GET", "/", nil)
	req.Host = host
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != expectedStatus {
		t.Errorf("%s - expected status %d, got %d", host, expectedStatus, rec.Code)
	}
}

func TestHostHeaderValidator(t *testing.T) {
	t.Run("valid Host header", func(t *testing.T) {
		testHostHeader(t, "yes.local", 200)
		testHostHeader(t, "Yes.local:8080", 200)
		testHostHeader(t, "127.0.0.1", 200)
		testHostHeader(t, "127.0.0.1:8080", 200)
	})
	t.Run("invalid Host header", func(t *testing.T) {
		testHostHeader(t, "no.local", 401)
		testHostHeader(t, "", 401)
	})
}
