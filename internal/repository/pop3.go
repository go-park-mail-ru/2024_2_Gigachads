package repository

import (
	"mail/config"
	"mail/pkg/pop3"
	"net/mail"
)

type POP3Repository struct {
	client *pop3.POP3Client
	cfg    *config.Config
}

func NewPOP3Repository(cfg *config.Config) *POP3Repository {
	return &POP3Repository{cfg: cfg}
}

func (p *POP3Repository) RetrieveMessages(from string) ([]*mail.Message, error) {
	return p.client.RetrieveMessages(from)
}
