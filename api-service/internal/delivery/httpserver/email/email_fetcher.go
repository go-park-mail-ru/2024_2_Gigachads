package email

import (
	"log/slog"
	"mail/api-service/internal/models"
	"time"
)

var NewTicker = time.NewTicker

type emailFetcher struct {
	emailService models.EmailUseCase
}

func NewEmailFetcher(es models.EmailUseCase) *emailFetcher {
	return &emailFetcher{
		emailService: es,
	}
}

func (ef *emailFetcher) Start() {
	if ef.emailService == nil {
		return
	}

	go func() {
		ticker := NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := ef.emailService.FetchEmailsViaPOP3(); err != nil {
			} else {
				slog.Info("Письма успешно получены через POP3")
			}
		}
	}()
}
