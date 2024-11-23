package usecase

import (
	"context"
	"fmt"
	"net/http"
	models "mail/api-service/internal/models"
	"mail/api-service/pkg/utils"
	"os"	
)

type AuthService struct {
	ms auth.AuthClient
}

func NewAuthService(authClient auth.AuthManagerClient) models.AuthUseCase {
	return &AuthService{
		ms: authClient,
	}
}

func (as *AuthService) Signup(ctx context.Context, signup *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	user, session, csrf, err := as.ms.Signup(ctx, signup)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (as *AuthService) Login(ctx context.Context, login *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	user, session, csrf, err := as.ms.Login(ctx, login)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, session, csrf, nil
}

func (as *UserService) Logout(ctx context.Context, email string) error {
	err := as.ms.Logout(ctx, email)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) CheckAuth(ctx context.Context, id string) (string, error) {
	email, err := as.ms.CheckAuth(ctx, id)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (as *AuthService) CheckCsrf(ctx context.Context, email string, csrf string) error {
	err := as.ms.CheckCsrf(ctx, email, csrf)
	if err != nil {
		return err
	}
	return nil
}

