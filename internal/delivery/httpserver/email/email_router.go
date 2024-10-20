package email

import (
	"github.com/gorilla/mux"
	"mail/internal/models"
)

type EmailRouter struct {
	EmailUseCase models.EmailUseCase
}

func NewEmailRouter(eu models.EmailUseCase) *EmailRouter {
	return &EmailRouter{EmailUseCase: eu}
}

func (er *EmailRouter) ConfigureEmailRouter(privateMux *mux.Router) {
	privateMux.HandleFunc("/mail/inbox", er.InboxHandler).Methods("GET", "OPTIONS")
}
