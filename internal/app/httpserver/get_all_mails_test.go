package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: после появления ручки создания письма в базе убрать
var testUserID = "test-uuid"

func TestHealthCheckHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllMails)

	req, err := http.NewRequest("GET", "/get_mails_by_user", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), "user-id", testUserID)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected, err := json.Marshal(testMockedItem)
	if err != nil {
		t.Errorf("cannot convert expected to json: %v", err)
	}
	if bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
