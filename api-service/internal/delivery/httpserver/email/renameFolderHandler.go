package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"net/http"
	"mail/api-service/internal/models"
	"github.com/gorilla/mux"
)


func (er *EmailRouter) RenameFolderHandler(w http.ResponseWriter, r *http.Request) {
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

	var folder models.RenameFolder
	err := json.NewDecoder(r.Body).Decode(&folder)
	
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	folderName = utils.Sanitize(folderName)
	folder.NewName = utils.Sanitize(folder.NewName)

	err = er.EmailUseCase.RenameFolder(email, folderName, folder.NewName)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
