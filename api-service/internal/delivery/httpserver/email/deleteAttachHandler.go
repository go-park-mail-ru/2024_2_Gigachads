package email

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
	"context"
	"time"
)

func (er *EmailRouter) DeleteAttachHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	err = er.EmailUseCase.DeleteAttach(ctx, file.Path)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_delete_file")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
