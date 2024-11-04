package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"strconv"
)

type DeleteEmailsRequest struct {
	IDs    []string `json:"ids"`
	Folder string   `json:"folder"`
}

func (er *EmailRouter) DeleteEmailsHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	userEmail := ctxEmail.(string)

	var request DeleteEmailsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "неверный формат данных")
		return
	}

	if len(request.IDs) == 0 {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "список ID пуст")
		return
	}

	if request.Folder == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "не указана папка")
		return
	}

	ids := make([]int, 0, len(request.IDs))
	for _, strID := range request.IDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "неверный формат ID")
			return
		}
		ids = append(ids, id)
	}

	if err := er.EmailUseCase.DeleteEmails(userEmail, ids, request.Folder); err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "ошибка при удалении писем")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
