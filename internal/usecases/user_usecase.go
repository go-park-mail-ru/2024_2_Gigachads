package usecase

import (
	models "mail/internal/models"
	repository "mail/internal/repository"
)

type UserUseCase struct {
	repo *repository.UserRepository
}

func NewUserUseCase(repo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uuc *UserUseCase) CreateUser(user *models.Signup) (*models.User, error) {
	return uuc.repo.CreateUser(user)
}

func (uuc *UserUseCase) GetUser(login *models.Login) (*models.User, error) {
	return uuc.repo.GetUser(login)
}
