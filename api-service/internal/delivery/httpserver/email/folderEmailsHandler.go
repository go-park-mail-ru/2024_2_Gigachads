package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"net/http"
	"mail/models"
)

func (er *EmailRouter) FolderEmailsHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	var folder models.Folder
	err := json.NewDecoder(r.Body).Decode(&folder)
	
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	mails, err := er.EmailUseCase.GetFolderEmails(email, folder.Name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "email_not_found")
		return
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