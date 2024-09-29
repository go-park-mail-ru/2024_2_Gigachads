package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testUserID = "test-uuid"

func TestGetAllMails(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllMails)
	req, err := http.NewRequest("GET", "/mail/inbox", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "user_id", Value: testUserID})
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var unmarshaledMails Mails
	err = json.Unmarshal(rr.Body.Bytes(), &unmarshaledMails)
	if err != nil {
		t.Errorf("cannot convert response body to struct: %v", err)
		return
	}
	if !unmarshaledMails.compare(mockedMails) {
		t.Errorf("handler returned unexpected body: got %v want %v", unmarshaledMails, mockedMails)
	}
}

func TestGetAllMails_Error(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllMails)
	req, err := http.NewRequest("GET", "/mail/inbox", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
	expected := errorResponse{
		Status: http.StatusForbidden,
		Body:   "Validation_error",
	}
	var result errorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("cannot convert response body to errorResponse struct: %v", err)
		return
	}
	if result != expected {
		t.Errorf("handler returned unexpected error response: got %v want %v", result, expected)
	}
}
