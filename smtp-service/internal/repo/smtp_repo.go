package repo

import (
	"mail/api-service/pkg/smtp"
	"mail/config"
)

type SMTPRepository struct {
	client *smtp.SMTPClient
	cfg    *config.Config
}

func NewSMTPRepository(client *smtp.SMTPClient, cfg *config.Config) *SMTPRepository {
	return &SMTPRepository{client: client, cfg: cfg}
}

func (s *SMTPRepository) SendEmail(from string, to []string, subject string, body string) error {
	return s.client.SendEmail(from, to, subject, body)
}
