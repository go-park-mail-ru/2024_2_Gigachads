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
	mux.HandleFunc("/email", er.SendEmailHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/email/{id}", er.SingleEmailHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/email/{id}/status", er.EmailStatusHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/email", er.DeleteEmailsHandler).Methods("DELETE", "OPTIONS")
	mux.HandleFunc("/folder", er.FoldersHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/folder/{name}", er.FolderEmailsHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/folder/{name}", er.RenameFolderHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/folder", er.CreateFolderHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/folder", er.DeleteFolderHandler).Methods("DELETE", "OPTIONS")
	mux.HandleFunc("/email/{id}/folder", er.ChangeEmailFolderHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/draft", er.CreateDraftHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/draft", er.UpdateDraftHandler).Methods("PUT", "OPTIONS")
	mux.HandleFunc("/draft/send", er.SendDraftHandler).Methods("POST", "OPTIONS")
}
