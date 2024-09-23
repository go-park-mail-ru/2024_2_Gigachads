package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HelloHandler(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", res.StatusCode)
	}

	expected := "hello, world!\n"
	if w.Body.String() != expected {
		t.Errorf("Expected response body '%s', got '%s'", expected, w.Body.String())
	}
}
