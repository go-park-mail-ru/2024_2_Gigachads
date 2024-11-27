package repo

import (
	"mail/api-service/pkg/logger"
	"mail/config"
	"mail/smtp-service/pkg/smtp"
)

type SMTPRepository struct {
	client *smtp.SMTPClient
	cfg    *config.Config
	logger logger.Logable
}

func NewSMTPRepository(client *smtp.SMTPClient, cfg *config.Config, l logger.Logable) *SMTPRepository {
	return &SMTPRepository{client: client, cfg: cfg, logger: l}
}

func (s *SMTPRepository) SendEmail(from string, to []string, subject string, body string) error {
	err := s.client.SendEmail(from, to, subject, body)
	if err != nil {
		s.logger.Error(err.Error())
	}
	return err
}
