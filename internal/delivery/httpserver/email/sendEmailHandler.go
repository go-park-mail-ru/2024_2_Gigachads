package email

import (
	"encoding/json"
	"mail/internal/delivery/converters"
	"mail/pkg/utils"
	"net/http"
)

func (er *EmailRouter) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	from := ctxEmail.(string)

	var req converters.SendEmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = er.EmailUseCase.SendEmail(from, []string{req.To}, req.Title, req.Body)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
