package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_EmailStatusHandler(t *testing.T) {
	tests := []struct {
		name       string
		setupAuth  bool
		emailID    string
		status     Status
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
		rawInput   string
	}{
		{
			name:      "успешное изменение статуса",
			setupAuth: true,
			emailID:   "1",
			status: Status{
				Status: true,
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					ChangeStatus(1, true).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "неавторизованный запрос",
			setupAuth:  false,
			emailID:    "1",
			wantStatus: http.StatusUnauthorized,
			wantBody: models.Error{
				Status: http.StatusUnauthorized,
				Body:   "unauthorized",
			},
		},
		{
			name:       "некорректный ID письма",
			setupAuth:  true,
			emailID:    "invalid",
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_path",
			},
		},
		{
			name:       "отсутствующий ID в пути",
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_path",
			},
		},
		{
			name:       "некорректный JSON в запросе",
			setupAuth:  true,
			emailID:    "1",
			rawInput:   "{invalid json",
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_json",
			},
		},
		{
			name:      "ошибка при изменении статуса",
			setupAuth: true,
			emailID:   "1",
			status: Status{
				Status: true,
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					ChangeStatus(1, true).
					Return(errors.New("status change error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_status",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockEmailUseCase)
			}

			router := NewEmailRouter(mockEmailUseCase)

			var reqBody []byte
			var err error
			if tt.rawInput != "" {
				reqBody = []byte(tt.rawInput)
			} else {
				reqBody, err = json.Marshal(tt.status)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/status/"+tt.emailID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Настройка маршрутизации с параметрами
			vars := map[string]string{}
			if tt.emailID != "" {
				vars["id"] = tt.emailID
			}
			req = mux.SetURLVars(req, vars)

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.EmailStatusHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus != http.StatusOK {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, errResponse)
			} else {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}
		})
	}
}
