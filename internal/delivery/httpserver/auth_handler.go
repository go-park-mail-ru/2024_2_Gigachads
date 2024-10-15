package httpserver

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	models "mail/internal/models"
	usecases "mail/internal/usecases"
	"net/http"
	"net/mail"
)

type AuthHandler struct {
	UserUseCase    *usecases.UserUseCase
	SessionUseCase *usecases.SessionUseCase
}

func NewAuthHandler(uu *usecases.UserUseCase, su *usecases.SessionUseCase) *AuthHandler {
	return &AuthHandler{UserUseCase: uu, SessionUseCase: su}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login models.Login

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	if !emailIsValid(login.Email) {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_input")
		return
	}

	if !inputIsValid(login.Password) {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_password")
		return
	}

	user, err := ah.UserUseCase.GetUser(&login)
	if err != nil {
		ErrorResponse(w, r, http.StatusForbidden, err.Error())
		return
	}

	session, err := ah.SessionUseCase.CreateSession(user.Email)

	cookie := http.Cookie{
		Name:     session.Name,
		Value:    session.ID,
		Expires:  session.Time,
		HttpOnly: true,
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signup models.Signup

	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	if !emailIsValid(signup.Email) {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_email")
		return
	}

	if !inputIsValid(signup.Name) || !inputIsValid(signup.Password) || !inputIsValid(signup.RePassword) {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_input")
		return
	}

	if signup.Password != signup.RePassword {
		ErrorResponse(w, r, http.StatusBadRequest, "invalid_password")
		return
	}

	user, err := ah.UserUseCase.CreateUser(&signup)

	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	session, err := ah.SessionUseCase.CreateSession(user.Email)

	cookie := http.Cookie{
		Name:     session.Name,
		Value:    session.ID,
		Expires:  session.Time,
		HttpOnly: true,
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	err = ah.SessionUseCase.DeleteSession(cookie.Value)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   cookie.Name,
		Value:  "",
		MaxAge: -1,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func emailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func inputIsValid(str string) bool {
	//match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", str)
	return true
}

func GenerateHash() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
