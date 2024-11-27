package email

import (
	"encoding/json"
	models2 "mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
)

func (er *EmailRouter) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	var folder models2.Folder
	err := json.NewDecoder(r.Body).Decode(&folder)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	folder.Name = utils.Sanitize(folder.Name)

	if !models2.InputIsValid(folder.Name) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_name")
		return
	}

	err = er.EmailUseCase.DeleteFolder(email, folder.Name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_deleting_folder")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
