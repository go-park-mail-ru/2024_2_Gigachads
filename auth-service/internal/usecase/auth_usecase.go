package usecase

import (
	"context"
	"fmt"
	"mail/auth-service/internal/models"
	models2 "mail/auth-service/internal/models"
	proto "mail/gen/go/auth"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	UserRepo    models.UserRepository
	SessionRepo models2.SessionRepository
	CsrfRepo    models2.CsrfRepository
}

func NewAuthServer(urepo models.UserRepository, srepo models2.SessionRepository, crepo models2.CsrfRepository) proto.AuthServiceServer {
	return &AuthServer{
		UserRepo:    urepo,
		SessionRepo: srepo,
		CsrfRepo:    crepo,
	}
}

func (as *AuthServer) Signup(ctx context.Context, signup *proto.SignupRequest) (*proto.SignupReply, error) {
	taken, err := as.UserRepo.IsExist(signup.GetEmail())
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, fmt.Errorf("login_taken")
	}
	user := &models.User{Name: signup.GetName(), Email: signup.GetEmail(), Password: signup.GetPassword()}
	userChecked, err := as.UserRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	session, err := as.SessionRepo.CreateSession(ctx, userChecked.Email)
	if err != nil {
		return nil, err
	}
	csrf, err := as.CsrfRepo.CreateCsrf(ctx, userChecked.Email)
	if err != nil {
		return nil, err
	}
	res := &proto.SignupReply{Session: session.ID, Csrf: csrf.ID}
	return res, nil
}

func (as *AuthServer) Login(ctx context.Context, login *proto.LoginRequest) (*proto.LoginReply, error) {
	taken, err := as.UserRepo.IsExist(login.GetEmail())
	if err != nil {
		return nil, err
	}
	if !taken {
		return nil, fmt.Errorf("user_does_not_exist")
	}
	user := &models.User{Email: login.GetEmail(), Password: login.GetPassword()}
	userChecked, err := as.UserRepo.CheckUser(user)
	if err != nil {
		return nil, err
	}
	session, err := as.SessionRepo.CreateSession(ctx, userChecked.Email)
	if err != nil {
		return nil, err
	}
	csrf, err := as.CsrfRepo.CreateCsrf(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	res := proto.LoginReply{Avatar: user.AvatarURL, SessionId: session.ID, CsrfId: csrf.ID, Name: user.Name}
	return &res, nil
}

func (as *AuthServer) Logout(ctx context.Context, logout *proto.LogoutRequest) (*proto.LogoutReply, error) {
	err := as.SessionRepo.DeleteSession(ctx, logout.GetEmail())
	if err != nil {
		return nil, err
	}
	err = as.CsrfRepo.DeleteCsrf(ctx, logout.GetEmail())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (as *AuthServer) CheckAuth(ctx context.Context, checkAuth *proto.AuthRequest) (*proto.AuthReply, error) {
	email, err := as.SessionRepo.GetSession(ctx, checkAuth.GetId())
	if err != nil {
		return nil, err
	}
	return &proto.AuthReply{Email: email}, nil
}

func (as *AuthServer) CheckCsrf(ctx context.Context, checkCsrf *proto.CsrfRequest) (*proto.CsrfReply, error) {
	email1, err := as.CsrfRepo.GetCsrf(ctx, checkCsrf.GetCsrf())
	if err != nil {
		return nil, err
	}
	email2, err := as.SessionRepo.GetSession(ctx, checkCsrf.GetEmail())
	if err != nil {
		return nil, err
	}
	if email1 != email2 {
		return nil, fmt.Errorf("invalid_csrf")
	}
	return nil, nil
}
