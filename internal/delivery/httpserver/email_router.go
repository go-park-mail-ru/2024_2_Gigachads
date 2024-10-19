package httpserver

import (
	"encoding/json"
	usecases "mail/internal/usecases"
	"net/http"
)

type EmailRouter struct {
	EmailUseCase usecases.EmailUseCase
}

func NewEmailRouter(eu usecases.EmailUseCase) *EmailRouter {
	return &EmailRouter{EmailUseCase: eu}
}

func (er *EmailRouter) InboxHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	mails, err := er.EmailUseCase.Inbox(cookie.Value)
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
