package auth

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/pkg/utils"
	"net/http"
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

	_, session, csrf, err := ar.UserUseCase.Signup(r.Context(), &signup)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	sessionCookie := http.Cookie{
		Name:     session.Name,
		Value:    session.ID,
		Expires:  session.Time,
		HttpOnly: true,
	}

	csrfCookie := http.Cookie{
		Name:     csrf.Name,
		Value:    csrf.ID,
		Expires:  csrf.Time,
		HttpOnly: true,
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &sessionCookie)
	http.SetCookie(w, &csrfCookie)
	w.WriteHeader(http.StatusOK)
}
