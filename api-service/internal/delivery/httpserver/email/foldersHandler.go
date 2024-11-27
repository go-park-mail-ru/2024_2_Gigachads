package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"net/http"
)

func (er *EmailRouter) FoldersHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	email := ctxEmail.(string)

	folders, err := er.EmailUseCase.GetFolders(email)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_getting_folders")
		return
	}

	result, err := json.Marshal(folders)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "json error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
