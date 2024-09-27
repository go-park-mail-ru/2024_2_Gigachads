package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type signuptodo struct {
	login		string	`json:"login"`
	name		string	`json:"name"`
	password	string	`json:"password"`
	repassword	string	`json:"repassword"`
}

func TestSignUpOK(t *testing.T) {

	todo := signuptodo{
	login: "aaaa",
	name: "bbbb",
	password: "cccc",
	repassword: "cccc",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignUpHandler)

	jsonReq, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestSignUpFailLogin(t *testing.T) {

	UserDB["nick@giga-mail.ru"] = User{login: "nick@giga-mail.ru", name: "nick", password: "12345"} //удалить, когда будет бд

	todo := signuptodo{
	login: "nick@giga-mail.ru",
	name: "bbbb",
	password: "cccc",
	repassword: "cccc",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignUpHandler)

	jsonReq, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

}

func TestSignUpFailPassword(t *testing.T) {

	todo := signuptodo{
	login: "aaaa",
	name: "bbbb",
	password: "cccc",
	repassword: "xxxx",
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignUpHandler)

	jsonReq, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

}