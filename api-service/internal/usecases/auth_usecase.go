package usecase

import (
	"context"
	"mail/api-service/internal/models"
	proto "mail/gen/go/auth"
)

type AuthService struct {
	ms proto.AuthServiceClient
}

func NewAuthService(authClient proto.AuthServiceClient) models.AuthUseCase {
	return &AuthService{
		ms: authClient,
	}
}

func (as *AuthService) Signup(ctx context.Context, signup *models.User) (string, string, error) {
	signupReq := &proto.SignupRequest{Name: signup.Name, Email: signup.Email, Password: signup.Password}
	signupReply, err := as.ms.Signup(ctx, signupReq)
	if err != nil {
		return "", "", err
	}
	return signupReply.GetSession(), signupReply.GetCsrf(), nil
}

func (as *AuthService) Login(ctx context.Context, login *models.User) (string, string, string, string, error) {
	loginReq := &proto.LoginRequest{Email: login.Email, Password: login.Password}
	loginReply, err := as.ms.Login(ctx, loginReq)
	if err != nil {
		return "", "", "", "", err
	}
	return loginReply.GetAvatar(), loginReply.GetName(), loginReply.GetSessionId(), loginReply.GetCsrfId(), nil
}

func (as *AuthService) Logout(ctx context.Context, email string) error {
	logoutReq := &proto.LogoutRequest{Email: email}
	_, err := as.ms.Logout(ctx, logoutReq)
	return err
}

func (as *AuthService) CheckAuth(ctx context.Context, id string) (string, error) {
	authReq := &proto.AuthRequest{Id: id}
	authReply, err := as.ms.CheckAuth(ctx, authReq)
	if err != nil {
		return "", err
	}
	return authReply.GetEmail(), nil
}

func (as *AuthService) CheckCsrf(ctx context.Context, session string, csrf string) error {
	csrfReq := &proto.CsrfRequest{Csrf: csrf, Email: session}
	_, err := as.ms.CheckCsrf(ctx, csrfReq)
	return err
}
