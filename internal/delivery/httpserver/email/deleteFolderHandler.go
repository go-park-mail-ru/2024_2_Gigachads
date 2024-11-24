package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"mail/internal/models"
	//"github.com/gorilla/mux"
)

func (er *EmailRouter) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) {
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

	folder.Name = utils.Sanitize(folder.Name)
		
	if !models.InputIsValid(folder.Name) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_name")
		return
	}

	err = er.EmailUseCase.DeleteFolder(email, folder.Name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
