package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type logintodo struct {
	login		string	`json:"login"`
	password	string	`json:"password"`
}

func TestLogInOK(t *testing.T) {

	UserDB["nick@giga-mail.ru"] = User{login: "nick@giga-mail.ru", name: "nick", password: "12345"} //убрать, когда будет бд

	todo := logintodo{
	login: "nick@giga-mail.ru",
	password: "12345",
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

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestLogInFailLogin(t *testing.T) {

	todo := logintodo{
	login: "aaaa",
	password: "cccc",
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

	//UserDB["nick@giga-mail.ru"] = User{login: "nick@giga-mail.ru", name: "nick", password: "12345"} //убрать, когда будет бд

	todo := logintodo{
	login: "nick@giga-mail.ru",
	password: "wrong",
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