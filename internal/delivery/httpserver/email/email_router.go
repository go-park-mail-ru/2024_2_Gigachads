package email

import (
	"mail/internal/models"

	"github.com/gorilla/mux"
)

type EmailRouter struct {
	EmailUseCase models.EmailUseCase
}

func NewEmailRouter(eu models.EmailUseCase) *EmailRouter {
	return &EmailRouter{EmailUseCase: eu}
}

func (er *EmailRouter) ConfigureEmailRouter(mux *mux.Router) {
	mux.HandleFunc("/email/inbox", er.InboxHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/email/sent", er.SentEmailsHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/email", er.SendEmailHandler).Methods("POST")
	mux.HandleFunc("/email/{id}", er.SingleEmailHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/email/{id}/status", er.EmailStatusHandler).Methods("PUT", "OPTIONS")
}
