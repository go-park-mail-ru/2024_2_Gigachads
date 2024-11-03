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

func (er *EmailRouter) ConfigureEmailRouter(privateMux *mux.Router) {
	privateMux.HandleFunc("/email/inbox", er.InboxHandler).Methods("GET", "OPTIONS")
	privateMux.HandleFunc("/email/sent", er.SentEmailsHandler).Methods("GET", "OPTIONS")
	privateMux.HandleFunc("/email", er.SendEmailHandler).Methods("POST")
	privateMux.HandleFunc("/email/{id}", er.SingleEmailHandler).Methods("GET", "OPTIONS")
	privateMux.HandleFunc("/email/{id}/status", er.EmailStatusHandler).Methods("PUT")
}
