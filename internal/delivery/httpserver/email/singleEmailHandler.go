package email

import (
	"encoding/json"
	"mail/internal/models"
	"mail/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (er *EmailRouter) SingleEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	strid, ok := vars["id"]
	if !ok {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_path")
		return
	}

	id, err := strconv.Atoi(strid)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_path")
		return
	}

	mail, err := er.EmailUseCase.GetEmailByID(id)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, "email_not_found")
		return
	}

	mails := []models.Email{mail}

	if mail.ParentID != 0 {
		parentMail, err := er.EmailUseCase.GetEmailByID(mail.ParentID)
		if err == nil {
			mails = append(mails, parentMail)
		}
	}

	result, err := json.Marshal(mails)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "json_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
