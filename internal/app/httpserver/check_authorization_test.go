package httpserver

import (
	"mail/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckAuthorizationByID_Success(t *testing.T) {
	testUser := database.User{
		Id:       "test-uuid",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	database.UserDB["test-uuid"] = testUser
	database.UserID["test-uuid"] = testUser.Name

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(checkAuthorizationByID)

	req, err := http.NewRequest("GET", "/auth/check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "test-uuid"})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Authorization successful\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCheckAuthorizationByID_NotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(checkAuthorizationByID)
	req, err := http.NewRequest("GET", "/auth/check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "nonexistent-uuid"})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}

	expected := "user not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
