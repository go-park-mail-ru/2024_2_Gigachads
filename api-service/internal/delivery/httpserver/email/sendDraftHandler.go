package email

import (
	"encoding/json"
	"mail/api-service/pkg/utils"
	"mail/models"
	"net/http"
	"strings"
	"time"
)

func (er *EmailRouter) SendDraftHandler(w http.ResponseWriter, r *http.Request) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	senderEmail := ctxEmail.(string)

	var req models.Draft
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "invalid_request_body")
		return
	}

	if req.ParentID == 0 {
		email := models.Email{
			ID:			  req.ID,
			Sender_email: senderEmail,
			Recipient:    req.Recipient,
			Title:        req.Title,
			Description:  req.Description,
			Sending_date: time.Now(),
			IsRead:       false,
			ParentID:     0,
		}

		err := er.EmailUseCase.SendDraft(email)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		er.EmailUseCase.SendEmail(
			r.Context(),
			senderEmail,
			[]string{req.Recipient},
			req.Title,
			req.Description,
		)
	} else {
		originalEmail, err := er.EmailUseCase.GetEmailByID(req.ParentID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, "parent_email_not_found")
			return
		}

		if strings.HasPrefix(req.Title, "Re:") {
			email := models.Email{
				ID:			  req.ID,
				Sender_email: senderEmail,
				Recipient:    originalEmail.Sender_email,
				Title:        req.Title,
				Description:  req.Description,
				Sending_date: time.Now(),
				IsRead:       false,
				ParentID:     req.ParentID,
			}

			err = er.EmailUseCase.SendDraft(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_reply")
				return
			}

			er.EmailUseCase.ReplyEmail(
				r.Context(),
				senderEmail,
				originalEmail.Sender_email,
				originalEmail,
				req.Description,
			)
		} else if strings.HasPrefix(req.Title, "Fwd:") {
			email := models.Email{
				ID:			  req.ID,
				Sender_email: senderEmail,
				Recipient:    req.Recipient,
				Title:        req.Title,
				Description:  req.Description,
				Sending_date: time.Now(),
				IsRead:       false,
				ParentID:     req.ParentID,
			}

			err = er.EmailUseCase.SendDraft(email)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusInternalServerError, "failed_to_save_forward")
				return
			}

			recipients := strings.Split(req.Recipient, ",")
			er.EmailUseCase.ForwardEmail(
				r.Context(),
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
