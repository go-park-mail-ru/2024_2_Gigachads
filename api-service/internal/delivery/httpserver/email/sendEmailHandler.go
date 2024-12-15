package email

import (
	"context"
	"encoding/json"
	"mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"net/http"
	"strings"
	"time"
)

func (er *EmailRouter) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
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
		// email := models.Email{
		// 	Sender_email: senderEmail,
		// 	Recipient:    req.Recipient,
		// 	Title:        req.Title,
		// 	Description:  req.Description,
		// 	Sending_date: time.Now(),
		// 	IsRead:       false,
		// 	ParentID:     0,
		// }
		email.Sender_email = senderEmail
		email.Sending_date = time.Now()
		email.IsRead = false

		err := er.EmailUseCase.SaveEmail(email)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		er.EmailUseCase.SendEmail(
			ctx,
			senderEmail,
			[]string{email.Recipient},
			email.Title,
			email.Description,
		)
	} else {
		originalEmail, err := er.EmailUseCase.GetEmailByID(email.ParentID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "parent_email_not_found")
			return
		}

		if strings.HasPrefix(email.Title, "Re:") {
			// email := models.Email{
			// 	Sender_email: senderEmail,
			// 	Recipient:    originalEmail.Sender_email,
			// 	Title:        req.Title,
			// 	Description:  req.Description,
			// 	Sending_date: time.Now(),
			// 	IsRead:       false,
			// 	ParentID:     req.ParentID,
			// }
			email.Sender_email = senderEmail
			email.Recipient = originalEmail.Sender_email
			email.Sending_date = time.Now()
			email.IsRead = false

			err = er.EmailUseCase.SaveEmail(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			er.EmailUseCase.ReplyEmail(
				ctx,
				senderEmail,
				originalEmail.Sender_email,
				originalEmail,
				email.Description,
			)
		} else if strings.HasPrefix(email.Title, "Fwd:") {
			// email := models.Email{
			// 	Sender_email: senderEmail,
			// 	Recipient:    req.Recipient,
			// 	Title:        req.Title,
			// 	Description:  req.Description,
			// 	Sending_date: time.Now(),
			// 	IsRead:       false,
			// 	ParentID:     req.ParentID,
			// }
			email.Sender_email = senderEmail
			email.Sending_date = time.Now()
			email.IsRead = false

			err = er.EmailUseCase.SaveEmail(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			recipients := strings.Split(email.Recipient, ",")
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			er.EmailUseCase.ForwardEmail(
				ctx,
				senderEmail,
				recipients,
				originalEmail,
			)
		} else {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_operation")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
