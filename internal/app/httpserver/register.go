package httpserver

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"mail/database"
	"net/http"
	"net/mail"
	"time"
	//"fmt"
)

type UserJSON struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	var user UserJSON

	// Декодируем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !emailIsValid(user.Email) {
		ErrorResponse(w, r, "invalid_email")
		return
	}

	if !inputIsValid(user.Name) {
		ErrorResponse(w, r, "invalid_input")
		return
	}

	if !inputIsValid(user.Password) {
		ErrorResponse(w, r, "invalid_input")
		return
	}

	if !inputIsValid(user.RePassword) {
		ErrorResponse(w, r, "invalid_input")
		return
	}

	if user.Password != user.RePassword {
		ErrorResponse(w, r, "invalid_password")
		return
	}

	if _, ok := database.UserDB[user.Email]; ok {
		ErrorResponse(w, r, "login_taken")
		return
	}

	uuid := GenerateUserID()
	database.UserDB[user.Email] = database.User{Id: uuid, Email: user.Email, Name: user.Name, Password: user.Password}
	database.UserID[uuid] = user.Email
	//w.Header().Set("Content-Type", "application/json")
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

func emailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func inputIsValid(str string) bool {
	// match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", str)
	return true
}

func GenerateUserID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
