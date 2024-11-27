package middleware

import (
	"context"
	"mail/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func TestLogMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		setupContext   func(context.Context) context.Context
		expectedStatus int
	}{
		{
			name:           "Normal request with existing requestID",
			method:         http.MethodGet,
			path:           "/test",
			setupContext:   func(ctx context.Context) context.Context { return context.WithValue(ctx, "requestID", "test-id") },
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Request without requestID",
			method:         http.MethodPost,
			path:           "/test",
			setupContext:   func(ctx context.Context) context.Context { return ctx },
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(MockLogger)
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()

			middleware := NewLogMW(mockLogger)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.expectedStatus)
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(tt.setupContext(req.Context()))
			rr := httptest.NewRecorder()

			handler := middleware.Handler(nextHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestPanicMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
	}{
		{
			name: "Handles panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			}),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Normal request without panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			handler := PanicMiddleware(tt.handler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestConfigureMWs(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTPServer.AllowedIPsByCORS = []string{"http://localhost:3000"}

	router := mux.NewRouter()
	mockAuth := new(MockAuthUseCase)
	authMW := NewAuthMW(mockAuth)

	handler := ConfigureMWs(cfg, router, authMW)
	assert.NotNil(t, handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, cfg.HTTPServer.AllowedIPsByCORS[0], rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestLogResponseWriter(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
	}{
		{
			name:         "Write response with status",
			statusCode:   http.StatusOK,
			responseBody: "test response",
		},
		{
			name:         "Write response with error status",
			statusCode:   http.StatusBadRequest,
			responseBody: "error response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			lrw := &LogResponseWriter{ResponseWriter: rr}

			lrw.WriteHeader(tt.statusCode)
			assert.Equal(t, tt.statusCode, lrw.statusCode)

			if tt.responseBody != "" {
				_, err := lrw.Write([]byte(tt.responseBody))
				assert.NoError(t, err)
				assert.Equal(t, tt.responseBody, rr.Body.String())
			}
		})
	}
}
