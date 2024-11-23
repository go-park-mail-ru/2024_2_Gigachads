package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"net/http"
)

func (er *EmailRouter) SentEmailsHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	emails, err := er.EmailUseCase.GetSentEmails(ctxEmail.(string))
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_get_sent_emails")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emails)
}
