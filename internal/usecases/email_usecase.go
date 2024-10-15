package usecase

import (
	models "mail/internal/models"
	repository "mail/internal/repository"
)

type EmailUseCase struct {
	repo *repository.EmailRepository
}

func NewEmailUseCase(repo *repository.EmailRepository) *EmailUseCase {
	return &EmailUseCase{
		repo: repo,
	}
}

func (euc *EmailUseCase) Inbox(id string) ([]models.Email, error) {
	return euc.repo.Inbox(id)
}
