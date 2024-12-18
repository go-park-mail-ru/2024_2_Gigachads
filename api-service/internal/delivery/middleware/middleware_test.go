package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "Успешная запись статус кода 200",
			statusCode: http.StatusOK,
		},
		{
			name:       "Успешная запись статус кода 404",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Успешная запись статус кода 500",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			lrw := &LogResponseWriter{
				ResponseWriter: rw,
				statusCode:     http.StatusOK,
			}

			lrw.WriteHeader(tt.statusCode)

			assert.Equal(t, tt.statusCode, lrw.statusCode)
			assert.Equal(t, tt.statusCode, rw.Code)
		})
	}
}

func TestNewLogMW(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)

	tests := []struct {
		name string
		log  logger.Logable
		want LogMiddleWare
	}{
		{
			name: "Успешное создание middleware логгера",
			log:  mockLogger,
			want: LogMiddleWare{logger: mockLogger},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogMW(tt.log)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLogMiddleWare_Handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)

	tests := []struct {
		name           string
		setupMock      func()
		request        *http.Request
		expectedStatus int
	}{
		{
			name: "Успешное логирование запроса без requestID",
			setupMock: func() {
				mockLogger.EXPECT().Info("User entered",
					"url", "/test",
					"method", "GET",
					"requestID", gomock.Any()).Times(1)
				mockLogger.EXPECT().Info("User left",
					"url", "/test",
					"status code", http.StatusOK,
					"requestID", gomock.Any()).Times(1)
			},
			request:        httptest.NewRequest("GET", "/test", nil),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Успешное логирование запроса с существующим requestID",
			setupMock: func() {
				mockLogger.EXPECT().Info("User entered",
					"url", "/test",
					"method", "POST",
					"requestID", "existing-id").Times(1)
				mockLogger.EXPECT().Info("User left",
					"url", "/test",
					"status code", http.StatusOK,
					"requestID", "existing-id").Times(1)
			},
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/test", nil)
				ctx := context.WithValue(req.Context(), "requestID", "existing-id")
				return req.WithContext(ctx)
			}(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			middleware := NewLogMW(mockLogger)
			handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем, что requestID был добавлен в контекст
				requestID := r.Context().Value("requestID")
				assert.NotNil(t, requestID)
				w.WriteHeader(tt.expectedStatus)
			}))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.request)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
