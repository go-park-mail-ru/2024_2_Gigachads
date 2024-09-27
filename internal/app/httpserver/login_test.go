package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogInOK(t *testing.T) {

	UserDB["nick@giga-mail.ru"] = User{ Name: "nick", Email: "nick@giga-mail.ru", Password: "12345"} //убрать, когда будет бд

	todo := UserLogin{
	Email: "nick@giga-mail.ru",
	Password: "12345",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogInHandler)

	jsonReq, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestLogInFailLogin(t *testing.T) {

	todo := UserLogin{
	Email: "vasia@giga-mail.ru",
	Password: "12345",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogInHandler)

	jsonReq, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

}

func TestLogInFailPassword(t *testing.T) {

	UserDB["nick@giga-mail.ru"] = User{ Name: "nick", Email: "nick@giga-mail.ru", Password: "12345"} //убрать, когда будет бд

	todo := UserLogin{
	Email: "nick@giga-mail.ru",
	Password: "12345678",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogInHandler)

	jsonReq, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

}