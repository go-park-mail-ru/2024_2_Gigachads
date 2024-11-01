package usecase

import (
	"fmt"
	"context"
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

func (us *UserService) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, error) {
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

	session, err := us.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

func (us *UserService) Login(ctx context.Context, login *models.User) (*models.User, *models.Session, error) {
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
	session, err := us.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

func (us *UserService) Logout(ctx context.Context, id string) error {
	err := us.SessionRepo.DeleteSession(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) CheckAuth(ctx context.Context, id string) (string, error) {
	session, err := us.SessionRepo.GetSession(ctx, id)
	if err != nil {
		return "", err
	}
	return session, nil
}
