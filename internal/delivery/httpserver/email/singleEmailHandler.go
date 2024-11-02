package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"mail/internal/models"
	"net/http"
	"strconv"
)

func (er *EmailRouter) SingleEmailHandler(w http.ResponseWriter, r *http.Request) {
	// ctxEmail := r.Context().Value("email")
	// if ctxEmail == nil {
	// 	utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	if !r.URL.Query().Has("id") {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_query")
		return
	}
	strid := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(strid)
	mail, _ := er.EmailUseCase.GetEmailByID(id)

	mails := make([]models.Email, 0)
	mails = append(mails, mail)
	for mail.ParentID != 0 {
		parent_id := mail.ParentID
		mail, _ := er.EmailUseCase.GetEmailByID(parent_id)
		mails = append(mails, mail)
	}

	result, err := json.Marshal(mails)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "json error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}