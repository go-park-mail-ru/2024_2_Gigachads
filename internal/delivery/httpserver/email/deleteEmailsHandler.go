package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
)

type DeleteEmailsRequest struct {
	IDs []int `json:"ids"`
}

func (er *EmailRouter) DeleteEmailsHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteEmailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "неверный формат данных")
		return
	}

	if len(req.IDs) == 0 {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "список ID пуст")
		return
	}

	if err := er.EmailUseCase.DeleteEmails(req.IDs); err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "ошибка при удалении писем")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
