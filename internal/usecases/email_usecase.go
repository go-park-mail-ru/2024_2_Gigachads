package usecase

import (
	"context"
	models "mail/internal/models"
)

type EmailService struct {
	EmailRepo   models.EmailRepository
	SessionRepo models.SessionRepository
}

func NewEmailService(erepo models.EmailRepository, srepo models.SessionRepository) models.EmailUseCase {
	return &EmailService{
		EmailRepo:   erepo,
		SessionRepo: srepo,
	}
}

func (es *EmailService) Inbox(ctx context.Context, sessionID string) ([]models.Email, error) {
	email, err := es.SessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return es.EmailRepo.Inbox(email)
}
