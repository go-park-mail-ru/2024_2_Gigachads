package user

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"mail/models"
	"net/http"
)

func (ar *UserRouter) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	var model *models.ChangePassword

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	model.Password = utils.Sanitize(model.Password)
	model.RePassword = utils.Sanitize(model.RePassword)
	model.OldPassword = utils.Sanitize(model.OldPassword)

	if !models.InputIsValid(model.Password) || !models.InputIsValid(model.RePassword) || !models.InputIsValid(model.OldPassword) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_input")
		return
	}
	if model.Password != model.RePassword {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_password")
		return
	}

	err = ar.UserUseCase.ChangePassword(email, model.Password)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
