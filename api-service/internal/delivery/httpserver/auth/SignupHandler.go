package auth

import (
	"context"
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
	"time"
)

func (ar *AuthRouter) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signup models.User

	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	signup.Email = utils.Sanitize(signup.Email)
	signup.Name = utils.Sanitize(signup.Name)
	signup.Password = utils.Sanitize(signup.Password)
	signup.RePassword = utils.Sanitize(signup.RePassword)

	if !models.EmailIsValid(signup.Email) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_email")
		return
	}

	if !models.InputIsValid(signup.Name) || !models.InputIsValid(signup.Password) || !models.InputIsValid(signup.RePassword) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_input")
		return
	}

	if signup.Password != signup.RePassword {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_password")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	sessionID, csrfID, err := ar.AuthUseCase.Signup(ctx, &signup)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_signup")
		return
	}

	sessionCookie := http.Cookie{
		Name:     "email",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    csrfID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &sessionCookie)
	http.SetCookie(w, &csrfCookie)
	w.WriteHeader(http.StatusOK)
}
