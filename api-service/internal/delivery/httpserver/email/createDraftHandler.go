package email

import (
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
	"strings"
	"time"
)

func (er *EmailRouter) CreateDraftHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	senderEmail := ctxEmail.(string)

	var email models.Email
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_request_body")
		return
	}

	if email.ParentID == 0 {
		email.Sender_email = senderEmail
		email.Sending_date = time.Now()
		email.IsRead = false

		err := er.EmailUseCase.CreateDraft(email)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusInternalServerError, "error_with_creating_draft")
			return
		}

	} else {
		originalEmail, err := er.EmailUseCase.GetEmailByID(email.ParentID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "parent_email_not_found")
			return
		}

		if strings.HasPrefix(email.Title, "Re:") {
			email.Sender_email = senderEmail
			email.Recipient = originalEmail.Sender_email
			email.Sending_date = time.Now()
			email.IsRead = false

			err = er.EmailUseCase.CreateDraft(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_reply")
				return
			}

		} else if strings.HasPrefix(email.Title, "Fwd:") {
			email.Sender_email = senderEmail
			email.Sending_date = time.Now()
			email.IsRead = false

			err = er.EmailUseCase.CreateDraft(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_forward")
				return
			}

		} else {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_operation")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
