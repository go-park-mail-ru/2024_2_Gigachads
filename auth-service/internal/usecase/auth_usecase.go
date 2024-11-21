package usecase

import (

	"context"
	"fmt"
	"net/http"
	models "mail/internal/models"
	"mail/pkg/utils"
	"os"	
)

type AuthServer struct {
	auth.UnimplementedAuthManagerServer
	UserRepo    models.UserRepository
	SessionRepo models.SessionRepository
	CsrfRepo models.CsrfRepository
}

func NewAuthServer(urepo models.UserRepository, srepo models.SessionRepository, crepo models.CsrfRepository) models.AuthServerUseCase {
	return &AuthServer{
		UserRepo: urepo,
		SessionRepo: srepo,
		CsrfRepo: crepo,
	}
}

func (as *AuthServer) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := as.UserRepo.IsExist(signup.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	if taken {
		return nil, nil, nil, fmt.Errorf("login_taken")
	}

	user, err := as.UserRepo.CreateUser(signup)
	if err != nil {
		return nil, nil, nil, err
	}

	session, err := as.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	csrf, err := as.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (as *AuthServer) Login(ctx context.Context, login *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	taken, err := as.UserRepo.IsExist(login.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	if !taken {
		return nil, nil, nil, fmt.Errorf("user_does_not_exist")
	}

	user, err := as.UserRepo.CheckUser(login)
	if err != nil {
		return nil, nil, nil, err
	}
	session, err := as.SessionRepo.CreateSession(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	csrf, err := as.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (as *UserService) Logout(ctx context.Context, email string) error {
	err := as.SessionRepo.DeleteSession(ctx, id)
	if err != nil {
		return err
	}
	err = as.CsrfRepo.DeleteCsrf(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthServer) CheckAuth(ctx context.Context, id string) (string, error) {
	session, err := as.SessionRepo.GetSession(ctx, id)
	if err != nil {
		return "", err
	}
	return session, nil
}

func (as *AuthServer) CheckCsrf(ctx context.Context, email string, csrf string) error {
	email1, err := as.CsrfRepo.GetCsrf(ctx, csrf)
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

