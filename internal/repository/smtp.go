package repository

import (
	"mail/config"
	"mail/pkg/smtp"
)

type SMTPRepository struct {
	client *smtp.SMTPClient
	cfg    *config.Config
}

func NewSMTPRepository() *SMTPRepository {
	return &SMTPRepository{}
}

func (s *SMTPRepository) SendEmail(from string, to []string, subject string, body string) error {
	return s.client.SendEmail(from, to, subject, body)
}
