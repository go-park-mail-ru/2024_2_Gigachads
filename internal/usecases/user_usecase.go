package usecase

import (
	"context"
	"fmt"
	models "mail/internal/models"
)

type UserService struct {
	UserRepo    models.UserRepository
	SessionRepo models.SessionRepository
	CsrfRepo    models.CsrfRepository
}

func NewUserService(urepo models.UserRepository, srepo models.SessionRepository, crepo models.CsrfRepository) models.UserUseCase {
	return &UserService{
		UserRepo:    urepo,
		SessionRepo: srepo,
		CsrfRepo:	 crepo,
	}
}

func (us *UserService) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := us.UserRepo.GetByEmail(signup.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	if taken {
		return nil, nil, nil, fmt.Errorf("login_taken")
	}

	user, err := us.UserRepo.CreateUser(signup)
	if err != nil {
		return nil, nil, nil, err
	}

	session, err := us.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	csrf, err := us.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (us *UserService) Login(ctx context.Context, login *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := us.UserRepo.GetByEmail(login.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	if !taken {
		return nil, nil, nil, fmt.Errorf("user_does_not_exist")
	}

	user, err := us.UserRepo.CheckUser(login)
	if err != nil {
		return nil, nil, nil, err
	}
	session, err := us.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	csrf, err := us.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (us *UserService) Logout(ctx context.Context, id string) error {
	err := us.SessionRepo.DeleteSession(ctx, id)
	if err != nil {
		return err
	}
	err = us.CsrfRepo.DeleteCsrf(ctx, id)
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

func (us *UserService) CheckCsrf(ctx context.Context, session string, csrf string) error {
	email1, err := us.CsrfRepo.GetCsrf(ctx, csrf)
	if err != nil {
		return err
	}
	email2, err := us.SessionRepo.GetSession(ctx, session)
	if err != nil {
		return err
	}
	if email1 != email2 {
		return fmt.Errorf("invalid_csrf")
	}
	return nil
}
