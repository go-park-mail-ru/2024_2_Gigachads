package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Status struct {
	Status bool `json:"status"`
}

func (er *EmailRouter) EmailStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

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

	var reqStatus Status
	err = json.NewDecoder(r.Body).Decode(&reqStatus)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_json")
		return
	}
	status := reqStatus.Status

	err = er.EmailUseCase.ChangeStatus(id, status)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_status")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
