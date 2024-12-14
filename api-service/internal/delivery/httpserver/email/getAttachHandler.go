package email

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
)

func (er *EmailRouter) GetAttachHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	var file models.FilePath
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	data, err := er.EmailUseCase.GetAttach(file.Path)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_get_file")
		return
	}

	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}