package repo

import (
	"mail/api-service/pkg/logger"
	"mail/config"
	"mail/models"
	"mail/smtp-service/pkg/pop3"
)

type POP3Repository struct {
	client *pop3.Pop3Client
	cfg    *config.Config
	logger logger.Logable
}

func NewPOP3Repository(client *pop3.Pop3Client, cfg *config.Config, l logger.Logable) *POP3Repository {
	return &POP3Repository{
		client: client,
		cfg:    cfg,
		logger: l,
	}
}

func (p *POP3Repository) Connect() error {
	err := p.client.Connect()
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}

func (p *POP3Repository) FetchEmails(repo models.EmailRepositorySMTP) error {
	err := p.client.FetchEmails(repo)
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}

func (p *POP3Repository) Quit() error {
	err := p.client.Quit()
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}
