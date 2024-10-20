package auth

import (
	"encoding/json"
	"mail/internal/models"
	"mail/pkg/utils"
	"net/http"
	"github.com/microcosm-cc/bluemonday"
)

func (ar *AuthRouter) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signup models.User
	sanitizer := bluemonday.UGCPolicy()

	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	signup.Email = sanitizer.Sanitize(signup.Email)
	signup.Name = sanitizer.Sanitize(signup.Name)
	signup.Password = sanitizer.Sanitize(signup.Password)
	signup.RePassword = sanitizer.Sanitize(signup.RePassword)
	
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

	_, session, err := ar.UserUseCase.Signup(&signup)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
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
