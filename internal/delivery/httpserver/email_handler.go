package httpserver

import (
	"encoding/json"
	usecases "mail/internal/usecases"
	"net/http"
)

type EmailHandler struct {
	EmailUseCase   *usecases.EmailUseCase
	SessionUseCase *usecases.SessionUseCase
	SMTPUsecase    *usecases.SMTPUsecase
}

func NewEmailHandler(eu *usecases.EmailUseCase, su *usecases.SessionUseCase) *EmailHandler {
	return &EmailHandler{EmailUseCase: eu, SessionUseCase: su}
}

func (eh *EmailHandler) Inbox(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	session, err := eh.SessionUseCase.GetSession(cookie.Value)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	mails, err := eh.EmailUseCase.Inbox(session.UserLogin)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := json.Marshal(mails)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, "json error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
