package email

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
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
	email := ctxEmail.(string)

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
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "email_not_found")
		return
	}

	if mail.Sender_email != email && mail.Recipient != email {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	mails := []models.Email{mail}

	for mail.ParentID != 0 {
		mail, err = er.EmailUseCase.GetEmailByID(mail.ParentID)
		if err == nil {
			mails = append(mails, mail)
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
