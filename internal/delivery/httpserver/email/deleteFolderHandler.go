package email

import (
	//"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"mail/internal/models"
	"github.com/gorilla/mux"
)

func (er *EmailRouter) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) {
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

	folderName = utils.Sanitize(folderName)
		
	if !models.InputIsValid(folderName) {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_name")
		return
	}

	err := er.EmailUseCase.DeleteFolder(email, folderName)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
