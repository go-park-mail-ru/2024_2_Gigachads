package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"net/http"
	"mail/api-service/internal/models"
)


func (er *EmailRouter) CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
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

	folder.Name = utils.Sanitize(folder.Name)
		
	if !models.InputIsValid(folder.Name) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_name")
		return
	}

	err = er.EmailUseCase.CreateFolder(email, folder.Name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}