package usecase

import (
	"context"
	"errors"
	"mail/api-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	proto "mail/gen/go/auth"
)

type mockAuthServiceClient struct {
	proto.AuthServiceClient
	signupResp *proto.SignupReply
	signupErr  error
	loginResp  *proto.LoginReply
	loginErr   error
	logoutErr  error
	authResp   *proto.AuthReply
	authErr    error
	csrfErr    error
}

func (m *mockAuthServiceClient) Signup(ctx context.Context, req *proto.SignupRequest, opts ...grpc.CallOption) (*proto.SignupReply, error) {
	return m.signupResp, m.signupErr
}

func (m *mockAuthServiceClient) Login(ctx context.Context, req *proto.LoginRequest, opts ...grpc.CallOption) (*proto.LoginReply, error) {
	return m.loginResp, m.loginErr
}

func (m *mockAuthServiceClient) Logout(ctx context.Context, req *proto.LogoutRequest, opts ...grpc.CallOption) (*proto.LogoutReply, error) {
	return &proto.LogoutReply{}, m.logoutErr
}

func (m *mockAuthServiceClient) CheckAuth(ctx context.Context, req *proto.AuthRequest, opts ...grpc.CallOption) (*proto.AuthReply, error) {
	return m.authResp, m.authErr
}

func (m *mockAuthServiceClient) CheckCsrf(ctx context.Context, req *proto.CsrfRequest, opts ...grpc.CallOption) (*proto.CsrfReply, error) {
	return &proto.CsrfReply{}, m.csrfErr
}

func TestNewAuthService(t *testing.T) {
	mockClient := &mockAuthServiceClient{}
	service := NewAuthService(mockClient)
	assert.NotNil(t, service, "Service should not be nil")
}

func TestSignup(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.User
		mockResp      *proto.SignupReply
		mockErr       error
		expectedSess  string
		expectedCsrf  string
		expectedError error
	}{
		{
			name: "Successful signup",
			user: &models.User{
				Email:    "test@test.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockResp: &proto.SignupReply{
				Session: "test-session",
				Csrf:    "test-csrf",
			},
			expectedSess: "test-session",
			expectedCsrf: "test-csrf",
		},
		{
			name: "Failed signup - email taken",
			user: &models.User{
				Email:    "existing@test.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockErr:       errors.New("email already exists"),
			expectedError: errors.New("email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAuthServiceClient{
				signupResp: tt.mockResp,
				signupErr:  tt.mockErr,
			}
			service := NewAuthService(mockClient)

			session, csrf, err := service.Signup(context.Background(), tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSess, session)
			assert.Equal(t, tt.expectedCsrf, csrf)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		user           *models.User
		mockResp       *proto.LoginReply
		mockErr        error
		expectedAvatar string
		expectedName   string
		expectedSess   string
		expectedCsrf   string
		expectedError  error
	}{
		{
			name: "Successful login",
			user: &models.User{
				Email:    "test@test.com",
				Password: "password123",
			},
			mockResp: &proto.LoginReply{
				Avatar:    "avatar.jpg",
				Name:      "Test User",
				SessionId: "test-session",
				CsrfId:    "test-csrf",
			},
			expectedAvatar: "avatar.jpg",
			expectedName:   "Test User",
			expectedSess:   "test-session",
			expectedCsrf:   "test-csrf",
		},
		{
			name: "Failed login - invalid credentials",
			user: &models.User{
				Email:    "test@test.com",
				Password: "wrongpass",
			},
			mockErr:       errors.New("invalid credentials"),
			expectedError: errors.New("invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAuthServiceClient{
				loginResp: tt.mockResp,
				loginErr:  tt.mockErr,
			}
			service := NewAuthService(mockClient)

			avatar, name, session, csrf, err := service.Login(context.Background(), tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAvatar, avatar)
			assert.Equal(t, tt.expectedName, name)
			assert.Equal(t, tt.expectedSess, session)
			assert.Equal(t, tt.expectedCsrf, csrf)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockErr       error
		expectedError error
	}{
		{
			name:  "Successful logout",
			email: "test@test.com",
		},
		{
			name:          "Failed logout",
			email:         "test@test.com",
			mockErr:       errors.New("logout failed"),
			expectedError: errors.New("logout failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAuthServiceClient{
				logoutErr: tt.mockErr,
			}
			service := NewAuthService(mockClient)

			err := service.Logout(context.Background(), tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCheckAuth(t *testing.T) {
	tests := []struct {
		name          string
		sessionID     string
		mockResp      *proto.AuthReply
		mockErr       error
		expectedEmail string
		expectedError error
	}{
		{
			name:      "Valid session",
			sessionID: "valid-session",
			mockResp: &proto.AuthReply{
				Email: "test@test.com",
			},
			expectedEmail: "test@test.com",
		},
		{
			name:          "Invalid session",
			sessionID:     "invalid-session",
			mockErr:       errors.New("invalid session"),
			expectedError: errors.New("invalid session"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAuthServiceClient{
				authResp: tt.mockResp,
				authErr:  tt.mockErr,
			}
			service := NewAuthService(mockClient)

			email, err := service.CheckAuth(context.Background(), tt.sessionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedEmail, email)
		})
	}
}

func TestCheckCsrf(t *testing.T) {
	tests := []struct {
		name          string
		session       string
		csrf          string
		mockErr       error
		expectedError error
	}{
		{
			name:    "Valid CSRF",
			session: "valid-session",
			csrf:    "valid-csrf",
		},
		{
			name:          "Invalid CSRF",
			session:       "valid-session",
			csrf:          "invalid-csrf",
			mockErr:       errors.New("invalid csrf token"),
			expectedError: errors.New("invalid csrf token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAuthServiceClient{
				csrfErr: tt.mockErr,
			}
			service := NewAuthService(mockClient)

			err := service.CheckCsrf(context.Background(), tt.session, tt.csrf)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
