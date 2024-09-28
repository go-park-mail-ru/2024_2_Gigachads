package httpserver

/*
TODO:
-кука
*/

import (
	"encoding/json"
	"net/http"
	"time"
	"mail/database"
	//"fmt"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {

	//database.UserDB["nick@giga-mail.ru"] = User{ Name: "nick", Email: "nick@giga-mail.ru", Password: "12345"} //убрать, когда будет бд

	var user UserLogin

	// Декодируем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	inputLogin := user.Email
	inputPassword := user.Password

	_, ok := database.UserDB[inputLogin]

	if !emailIsValid(inputLogin) {
		ErrorResponse(w, r, "invalid_input")
		return
	}

	if !inputIsValid(inputPassword) {
		ErrorResponse(w, r, "invalid_input")
		return
	}

	if !ok {
		ErrorResponse(w, r, "user_does_not_exist")
		return
	}

	if database.UserDB[inputLogin].Password != inputPassword {
		ErrorResponse(w, r, "invalid_password")
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "user_id",
		Value:    database.UserDB[user.Email].Id,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
