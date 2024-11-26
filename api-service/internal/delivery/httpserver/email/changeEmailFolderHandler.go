package email

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mail/api-service/pkg/utils"
	"mail/models"
	"net/http"
	"strconv"
)

func (er *EmailRouter) ChangeEmailFolderHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	vars := mux.Vars(r)
	strid, ok := vars["id"]
	if !ok {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_path")
		return
	}
	id, err := strconv.Atoi(strid)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_path")
		return
	}

	var folder models.Folder
	err = json.NewDecoder(r.Body).Decode(&folder)
	
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	folder.Name = utils.Sanitize(folder.Name)

	err = er.EmailUseCase.ChangeEmailFolder(id, email, folder.Name)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
