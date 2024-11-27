package usecase

import (
	"context"
	"mail/auth-service/internal/mocks"
	"mail/auth-service/internal/models"
	"mail/gen/go/auth"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	csrfRepo := mocks.NewMockCsrfRepo(ctrl)

	server := NewAuthServer(userRepo, sessionRepo, csrfRepo)
	ctx := context.Background()

	t.Run("успешная регистрация", func(t *testing.T) {
		req := &auth.SignupRequest{
			Email:    "test@test.com",
			Password: "password",
			Name:     "Test User",
		}

		userRepo.EXPECT().IsExist(req.Email).Return(false, nil)
		userRepo.EXPECT().CreateUser(gomock.Any()).Return(&models.User{
			Email: req.Email,
			Name:  req.Name,
		}, nil)
		sessionRepo.EXPECT().CreateSession(ctx, req.Email).Return(&models.Session{ID: "session-id"}, nil)
		csrfRepo.EXPECT().CreateCsrf(ctx, req.Email).Return(&models.Csrf{ID: "csrf-id"}, nil)

		resp, err := server.Signup(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "session-id", resp.Session)
		assert.Equal(t, "csrf-id", resp.Csrf)
	})

	t.Run("email уже существует", func(t *testing.T) {
		req := &auth.SignupRequest{
			Email:    "exists@test.com",
			Password: "password",
			Name:     "Test User",
		}

		userRepo.EXPECT().IsExist(req.Email).Return(true, nil)

		resp, err := server.Signup(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "login_taken", err.Error())
		assert.Nil(t, resp)
	})
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	csrfRepo := mocks.NewMockCsrfRepo(ctrl)

	server := NewAuthServer(userRepo, sessionRepo, csrfRepo)
	ctx := context.Background()

	t.Run("успешный логаут", func(t *testing.T) {
		req := &auth.LogoutRequest{
			Email: "test@test.com",
		}

		sessionRepo.EXPECT().DeleteSession(ctx, req.Email).Return(nil)
		csrfRepo.EXPECT().DeleteCsrf(ctx, req.Email).Return(nil)

		_, err := server.Logout(ctx, req)
		assert.NoError(t, err)
	})
}

func TestCheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	csrfRepo := mocks.NewMockCsrfRepo(ctrl)

	server := NewAuthServer(userRepo, sessionRepo, csrfRepo)
	ctx := context.Background()

	t.Run("успешная проверка сессии", func(t *testing.T) {
		req := &auth.AuthRequest{
			Id: "session-id",
		}

		sessionRepo.EXPECT().GetSession(ctx, req.Id).Return("test@test.com", nil)

		resp, err := server.CheckAuth(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "test@test.com", resp.Email)
	})
}

func TestCheckCsrf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	csrfRepo := mocks.NewMockCsrfRepo(ctrl)

	server := NewAuthServer(userRepo, sessionRepo, csrfRepo)
	ctx := context.Background()

	t.Run("успешная проверка CSRF", func(t *testing.T) {
		req := &auth.CsrfRequest{
			Csrf:  "csrf-token",
			Email: "session-id",
		}

		sessionRepo.EXPECT().GetSession(ctx, req.Email).Return("test@test.com", nil)
		csrfRepo.EXPECT().GetCsrf(ctx, req.Csrf).Return("test@test.com", nil)

		_, err := server.CheckCsrf(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("несовпадающие email", func(t *testing.T) {
		req := &auth.CsrfRequest{
			Csrf:  "csrf-token",
			Email: "session-id",
		}

		sessionRepo.EXPECT().GetSession(ctx, req.Email).Return("test2@test.com", nil)
		csrfRepo.EXPECT().GetCsrf(ctx, req.Csrf).Return("test@test.com", nil)

		_, err := server.CheckCsrf(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "invalid_csrf", err.Error())
	})
}
