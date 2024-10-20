package auth

import (
	"encoding/json"
	"mail/internal/models"
	"mail/pkg/utils"
	"net/http"
)

func (ar *AuthRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login models.User

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	if !models.EmailIsValid(login.Email) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_input")
		return
	}

	if !models.InputIsValid(login.Password) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_password")
		return
	}

	_, session, err := ar.UserUseCase.Login(&login)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusForbidden, err.Error())
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
