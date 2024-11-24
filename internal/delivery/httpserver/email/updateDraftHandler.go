package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"mail/internal/models"
)


func (er *EmailRouter) UpdateDraftHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	var draft models.Draft
	err := json.NewDecoder(r.Body).Decode(&draft)
	
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}

	draft.Title = utils.Sanitize(draft.Title)
	draft.Description = utils.Sanitize(draft.Description)

	err = er.EmailUseCase.UpdateDraft(draft)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}