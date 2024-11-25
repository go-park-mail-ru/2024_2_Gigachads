package repo

import (
	"mail/api-service/pkg/pop3"
	"mail/config"
	"mail/models"
)

type POP3Repository struct {
	client *pop3.Pop3Client
	cfg    *config.Config
}

func NewPOP3Repository(client *pop3.Pop3Client, cfg *config.Config) *POP3Repository {
	return &POP3Repository{
		client: client,
		cfg:    cfg,
	}
}

func (p *POP3Repository) Connect() error {
	return p.client.Connect()
}

func (p *POP3Repository) FetchEmails(repo models.EmailRepository) error {
	return p.client.FetchEmails(repo)
}

func (p *POP3Repository) Quit() error {
	return p.client.Quit()
}
