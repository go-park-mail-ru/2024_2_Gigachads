package usecase

import (
	"fmt"
	models "mail/internal/models"
)

type UserService struct {
	UserRepo    models.UserRepository
	SessionRepo models.SessionRepository
}

func NewUserService(urepo models.UserRepository, srepo models.SessionRepository) models.UserUseCase {
	return &UserService{
		UserRepo:    urepo,
		SessionRepo: srepo,
	}
}

func (us *UserService) Signup(signup *models.User) (*models.User, *models.Session, error) {
	taken, err := us.UserRepo.GetByEmail(signup.Email)
	if err != nil {
		return nil, nil, err
	}
	if taken {
		return nil, nil, fmt.Errorf("login_taken")
	}

	user, err := us.UserRepo.CreateUser(signup)
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
	taken, err := us.UserRepo.GetByEmail(login.Email)
	if err != nil {
		return nil, nil, err
	}
	if !taken {
		return nil, nil, fmt.Errorf("user_does_not_exist")
	}

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
