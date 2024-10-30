package email

import (
	"encoding/json"
	"net/http"
	"strings"

	"mail/internal/delivery/converters"
)

func (er *EmailRouter) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	senderEmail := r.Context().Value("email")
	if senderEmail == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req converters.SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Если ParentId == 0, то это новое письмо
	if req.ParentId == 0 {
		err := er.EmailUseCase.SendEmail(
			senderEmail.(string),
			[]string{req.Recipient},
			req.Title,
			req.Description,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// Иначе определяем действием: ответ или пересылка
		originalEmail, err := er.EmailUseCase.GetEmailByID(req.ParentId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.HasPrefix(req.Title, "Re:") {
			// Ответ
			err = er.EmailUseCase.ReplyEmail(
				senderEmail.(string),
				originalEmail.Sender_email,
				originalEmail,
				req.Description,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if strings.HasPrefix(req.Title, "Fwd:") {
			// Пересылка
			recipients := strings.Split(req.Recipient, ",")
			err = er.EmailUseCase.ForwardEmail(
				senderEmail.(string),
				recipients,
				originalEmail,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}