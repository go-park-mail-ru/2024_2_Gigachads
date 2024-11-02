package email

import (
	"encoding/json"
	"mail/pkg/utils"
	//"mail/internal/models"
	"net/http"
	"strconv"
)

type Status struct{
	status string
}

func (er *EmailRouter) EmailStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	if !r.URL.Query().Has("id") {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_query")
		return
	}
	strid := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(strid)

	var reqStatus Status
	err := json.NewDecoder(r.Body).Decode(&reqStatus)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}
	status := reqStatus.status

	err = er.EmailUseCase.ChangeStatus(id, status)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_status")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}