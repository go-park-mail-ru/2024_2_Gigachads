package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUpOK(t *testing.T) {

	todo := UserJSON{
	Name: "aaaa",
	Email: "petia@giga-mail.ru",
	Password: "cccc",
	RePassword: "cccc",
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

	UserDB["nick@giga-mail.ru"] = User{ Name: "nick", Email: "nick@giga-mail.ru", Password: "12345"} //убрать, когда будет бд

	todo := UserJSON{
	Name: "aaaa",
	Email: "nick@giga-mail.ru",
	Password: "cccc",
	RePassword: "cccc",
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

	todo := UserJSON{
	Name: "aaaa",
	Email: "nick@giga-mail.ru",
	Password: "cccc",
	RePassword: "dddd",
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