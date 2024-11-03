package email

import (
	"encoding/json"
	"mail/internal/delivery/converters"
	"mail/internal/models"
	"mail/pkg/utils"
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

	var req converters.SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_request_body")
		return
	}

	if req.ParentId == 0 {
		email := models.Email{
			Sender_email: senderEmail,
			Recipient:    req.Recipient,
			Title:        req.Title,
			Description:  req.Description,
			Sending_date: time.Now(),
			IsRead:       false,
		}

		err := er.EmailUseCase.SaveEmail(email)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_email")
			return
		}

		err = er.EmailUseCase.SendEmail(
			senderEmail,
			[]string{req.Recipient},
			req.Title,
			req.Description,
		)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_send_email")
			return
		}
	} else {
		originalEmail, err := er.EmailUseCase.GetEmailByID(req.ParentId)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "parent_email_not_found")
			return
		}

		if strings.HasPrefix(req.Title, "Re:") {
			email := models.Email{
				Sender_email: senderEmail,
				Recipient:    originalEmail.Sender_email,
				Title:        req.Title,
				Description:  req.Description,
				Sending_date: time.Now(),
				IsRead:       false,
			}

			err = er.EmailUseCase.SaveEmail(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_reply")
				return
			}

			err = er.EmailUseCase.ReplyEmail(
				senderEmail,
				originalEmail.Sender_email,
				originalEmail,
				req.Description,
			)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_reply")
				return
			}
		} else if strings.HasPrefix(req.Title, "Fwd:") {
			email := models.Email{
				Sender_email: senderEmail,
				Recipient:    req.Recipient,
				Title:        req.Title,
				Description:  originalEmail.Description,
				Sending_date: time.Now(),
				IsRead:       false,
			}

			err = er.EmailUseCase.SaveEmail(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_forward")
				return
			}

			recipients := strings.Split(req.Recipient, ",")
			err = er.EmailUseCase.ForwardEmail(
				senderEmail,
				recipients,
				originalEmail,
			)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_forward")
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
