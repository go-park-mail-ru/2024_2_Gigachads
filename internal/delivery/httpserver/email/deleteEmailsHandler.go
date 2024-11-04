package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
	"strconv"
)

func (er *EmailRouter) DeleteEmailsHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	userEmail := ctxEmail.(string)

	var strIDs []string
	if err := json.NewDecoder(r.Body).Decode(&strIDs); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "неверный формат данных")
		return
	}

	if len(strIDs) == 0 {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "список ID пуст")
		return
	}

	ids := make([]int, 0, len(strIDs))
	for _, strID := range strIDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "неверный формат ID")
			return
		}
		ids = append(ids, id)
	}

	if err := er.EmailUseCase.DeleteEmails(userEmail, ids); err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "ошибка при удалении писем")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
