package httpserver

import (
	"encoding/json"
	models "mail/internal/models"
	usecases "mail/internal/usecases"
	"net/http"
	"net/mail"
)

type AuthRouter struct {
	UserUseCase usecases.UserUseCase
}

func NewAuthRouter(uu usecases.UserUseCase) *AuthRouter {
	return &AuthRouter{UserUseCase: uu}
}

func (ar *AuthRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login models.User

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

	_, session, err := ar.UserUseCase.Login(&login)
	if err != nil {
		ErrorResponse(w, r, http.StatusForbidden, err.Error())
		return
	}

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

func (ar *AuthRouter) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signup models.User

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

	_, session, err := ar.UserUseCase.Signup(&signup)

	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

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

func (ar *AuthRouter) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	err = ar.UserUseCase.Logout(cookie.Value)
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
