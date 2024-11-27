package email

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
)

func (er *EmailRouter) RenameFolderHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	var folder models.RenameFolder
	err := json.NewDecoder(r.Body).Decode(&folder)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	folder.Name = utils.Sanitize(folder.Name)
	folder.NewName = utils.Sanitize(folder.NewName)

	err = er.EmailUseCase.RenameFolder(email, folder.Name, folder.NewName)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_rename_folder")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
