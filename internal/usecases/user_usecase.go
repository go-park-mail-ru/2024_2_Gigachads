package usecase

import (
	models "mail/internal/models"
	repository "mail/internal/repository"
)

type UserUseCase interface {
	Signup(user *models.User) (*models.User, *models.Session, error)
	Login(user *models.User) (*models.User, *models.Session, error)
	Logout(id string) error
	CheckAuth(sessionID string) (*models.Session, error)
}

type UserService struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
}

func NewUserService(urepo repository.UserRepository, srepo repository.SessionRepository) UserUseCase {
	return &UserService{
		UserRepo:    urepo,
		SessionRepo: srepo,
	}
}

func (us *UserService) Signup(user *models.User) (*models.User, *models.Session, error) {
	user, err := us.UserRepo.CreateUser(user)
	if err != nil {
		return nil, nil, err
	}
	session, err := us.SessionRepo.CreateSession(user.Email)
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

func (us *UserService) Login(login *models.User) (*models.User, *models.Session, error) {
	user, err := us.UserRepo.CheckUser(login)
	if err != nil {
		return nil, nil, err
	}
	session, err := us.SessionRepo.CreateSession(user.Email)
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

func (us *UserService) Logout(id string) error {
	err := us.SessionRepo.DeleteSession(id)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) CheckAuth(id string) (*models.Session, error) {
	session, err := us.SessionRepo.GetSession(id)
	if err != nil {
		return nil, err
	}
	return session, nil
}
