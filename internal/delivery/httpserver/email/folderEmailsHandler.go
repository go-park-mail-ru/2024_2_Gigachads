package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"mail/internal/models"
	//"github.com/gorilla/mux"
)

func (er *EmailRouter) FolderEmailsHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	// vars := mux.Vars(r)
	// folderName, ok := vars["name"]
	// if !ok {
	// 	utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_path")
	// 	return
	// }

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