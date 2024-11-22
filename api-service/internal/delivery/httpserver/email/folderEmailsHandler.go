package email

import (
	"encoding/json"
	"mail/internal/models"
	"mail/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (er *EmailRouter) FolderEmailsHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	vars := mux.Vars(r)
	folderName, ok := vars["name"]
	if !ok {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_path")
		return
	}

	mails, err := er.EmailUseCase.GetFolderEmails(email, folderName)
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