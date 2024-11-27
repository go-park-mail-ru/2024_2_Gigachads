package delivery

import (
	"log/slog"
	"mail/smtp-service/internal/models"
	"time"
)

var NewTicker = time.NewTicker

type EmailFetcher struct {
	Pop3Repo  models.POP3Repository
	EmailRepo models.EmailRepositorySMTP
}

func NewEmailFetcher(prepo models.POP3Repository, erepo models.EmailRepositorySMTP) *EmailFetcher {
	return &EmailFetcher{
		Pop3Repo:  prepo,
		EmailRepo: erepo,
	}
}
func StartEmailFetcher(prepo models.POP3Repository, erepo models.EmailRepositorySMTP) {
	fetcher := NewEmailFetcher(prepo, erepo)
	fetcher.Start()
}

func (ef *EmailFetcher) Start() {
	if ef.Pop3Repo == nil || ef.EmailRepo == nil {
		return
	}

	go func() {
		ticker := NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := ef.FetchEmailsViaPOP3(); err != nil {
			} else {
				slog.Info("Письма успешно получены через POP3")
			}
		}
	}()
}
func (ef *EmailFetcher) FetchEmailsViaPOP3() error {
	if err := ef.Pop3Repo.Connect(); err != nil {
		return err
	}
	defer ef.Pop3Repo.Quit()

	if err := ef.Pop3Repo.FetchEmails(ef.EmailRepo); err != nil {
		return err
	}

	return nil
}
