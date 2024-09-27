package httpserver

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"regexp"
	"crypto/rand"
    "encoding/hex"
    "time"
	//"fmt"
)

type UserJSON struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var UserDB = make(map[string]User)

var UserID = make(map[string]string)

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

	if _, ok := UserDB[user.Email]; ok {
		ErrorResponse(w, r, "login_taken")
		return
	}

	UserDB[user.Email] = User{Email: user.Email, Name: user.Name, Password: user.Password}
	UserID[user.Email] = GenerateUserID()
	//w.Header().Set("Content-Type", "application/json")
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "user_id",
		Value:    UserID[user.Email],
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
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", str)
	return match
}

func GenerateUserID() string {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        panic(err)
    }
    return hex.EncodeToString(bytes)
}