package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"strconv"
)

type DeleteEmailsRequest struct {
	IDs []string `json:"ids"`
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

	ids := make([]int, 0, len(req.IDs))
	for _, strID := range req.IDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "неверный формат ID")
			return
		}
		ids = append(ids, id)
	}

	if err := er.EmailUseCase.DeleteEmails(ids); err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "ошибка при удалении писем")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
