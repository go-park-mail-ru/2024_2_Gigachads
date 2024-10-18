package usecase

import (
	"mail/internal/repository"
	"net/mail"
)

type POP3Usecase struct {
	pop3Repository *repository.POP3Repository
}

func NewPOP3Usecase(pop3Repository *repository.POP3Repository) *POP3Usecase {
	return &POP3Usecase{pop3Repository: pop3Repository}
}

func (p *POP3Usecase) RetrieveMessages(from string) ([]*mail.Message, error) {
	return p.pop3Repository.RetrieveMessages(from)
}
