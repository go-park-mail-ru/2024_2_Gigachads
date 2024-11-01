package email

import (
	"log/slog"
	"mail/internal/models"
	"time"
)

var newTicker = time.NewTicker

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
		ticker := newTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := ef.emailService.FetchEmailsViaPOP3(); err != nil {
				slog.Error("Ошибка при получении писем через POP3", "error", err)
			} else {
				slog.Info("Письма успешно получены через POP3")
			}
		}
	}()
}
