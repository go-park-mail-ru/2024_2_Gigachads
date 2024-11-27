package auth

import (
	"context"
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
	"time"
)

func (ar *AuthRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login models.User

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	login.Password = utils.Sanitize(login.Password)
	login.Email = utils.Sanitize(login.Email)

	if !models.EmailIsValid(login.Email) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_input")
		return
	}

	if !models.InputIsValid(login.Password) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_password")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	avatar, name, session, csrf, err := ar.AuthUseCase.Login(ctx, &login)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusForbidden, "invalid_login_or_password")
		return
	}

	userLogin := models.UserLogin{Email: login.Email, Name: name, AvatarURL: avatar}
	if userLogin.AvatarURL == "" {
		userLogin.AvatarURL = "/icons/default.png"
	}

	result, err := json.Marshal(userLogin)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "json_error")
		return
	}

	sessionCookie := http.Cookie{
		Name:     "email",
		Value:    session,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    csrf,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &sessionCookie)
	http.SetCookie(w, &csrfCookie)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
