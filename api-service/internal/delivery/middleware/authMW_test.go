package middleware

import (
	"context"
	"errors"
	"mail/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Signup(ctx context.Context, signup *models.User) (string, string, error) {
	args := m.Called(ctx, signup)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthUseCase) Login(ctx context.Context, login *models.User) (string, string, string, string, error) {
	args := m.Called(ctx, login)
	return args.String(0), args.String(1), args.String(2), args.String(3), args.Error(4)
}

func (m *MockAuthUseCase) Logout(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthUseCase) CheckAuth(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) CheckCsrf(ctx context.Context, session string, csrf string) error {
	args := m.Called(ctx, session, csrf)
	return args.Error(0)
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		cookies        map[string]string
		setupMock      func(*MockAuthUseCase)
		expectedEmail  string
		expectedStatus int
	}{
		{
			name:   "OPTIONS request bypasses auth",
			method: http.MethodOptions,
			cookies: map[string]string{
				"email": "test-session",
				"csrf":  "test-csrf",
			},
			setupMock: func(m *MockAuthUseCase) {},
		},
		{
			name:   "Successful auth",
			method: http.MethodGet,
			cookies: map[string]string{
				"email": "test-session",
				"csrf":  "test-csrf",
			},
			setupMock: func(m *MockAuthUseCase) {
				m.On("CheckAuth", mock.Anything, "test-session").Return("test@test.com", nil)
				m.On("CheckCsrf", mock.Anything, "test-session", "test-csrf").Return(nil)
			},
			expectedEmail: "test@test.com",
		},
		{
			name:   "Missing email cookie",
			method: http.MethodGet,
			cookies: map[string]string{
				"csrf": "test-csrf",
			},
			setupMock:      func(m *MockAuthUseCase) {},
			expectedEmail:  "",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Missing csrf cookie",
			method: http.MethodGet,
			cookies: map[string]string{
				"email": "test-session",
			},
			setupMock:      func(m *MockAuthUseCase) {},
			expectedEmail:  "",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid session",
			method: http.MethodGet,
			cookies: map[string]string{
				"email": "invalid-session",
				"csrf":  "test-csrf",
			},
			setupMock: func(m *MockAuthUseCase) {
				m.On("CheckAuth", mock.Anything, "invalid-session").Return("", errors.New("invalid session"))
			},
			expectedEmail:  "",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid csrf",
			method: http.MethodGet,
			cookies: map[string]string{
				"email": "test-session",
				"csrf":  "invalid-csrf",
			},
			setupMock: func(m *MockAuthUseCase) {
				m.On("CheckAuth", mock.Anything, "test-session").Return("test@test.com", nil)
				m.On("CheckCsrf", mock.Anything, "test-session", "invalid-csrf").Return(errors.New("invalid csrf"))
			},
			expectedEmail:  "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(MockAuthUseCase)
			tt.setupMock(mockAuth)

			middleware := NewAuthMW(mockAuth)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				email := r.Context().Value("email")
				if email != nil {
					assert.Equal(t, tt.expectedEmail, email.(string))
				} else {
					assert.Empty(t, tt.expectedEmail)
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, "/test", nil)
			for name, value := range tt.cookies {
				req.AddCookie(&http.Cookie{Name: name, Value: value})
			}

			rr := httptest.NewRecorder()

			handler := middleware.Handler(nextHandler)
			handler.ServeHTTP(rr, req)

			if tt.expectedStatus != 0 {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestNewAuthMW(t *testing.T) {
	mockAuth := new(MockAuthUseCase)
	middleware := NewAuthMW(mockAuth)
	assert.NotNil(t, middleware)
	assert.Equal(t, mockAuth, middleware.AuthUseCase)
}
