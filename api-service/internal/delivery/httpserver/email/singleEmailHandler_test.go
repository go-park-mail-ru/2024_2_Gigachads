package email

import (
	"context"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_SingleEmailHandler(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		setupAuth  bool
		emailID    string
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
	}{
		{
			name:      "успешное получение письма без родительских",
			setupAuth: true,
			emailID:   "1",
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{
						ID:           1,
						Sender_email: "test@example.com",
						Recipient:    "recipient@example.com",
						Title:        "Test Email",
						Description:  "Test Content",
						Sending_date: testTime,
						ParentID:     0,
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: []models.Email{
				{
					ID:           1,
					Sender_email: "test@example.com",
					Recipient:    "recipient@example.com",
					Title:        "Test Email",
					Description:  "Test Content",
					Sending_date: testTime,
					ParentID:     0,
				},
			},
		},
		{
			name:      "успешное получение цепочки писем",
			setupAuth: true,
			emailID:   "2",
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(2).
					Return(models.Email{
						ID:           2,
						Sender_email: "test@example.com",
						Recipient:    "recipient@example.com",
						Title:        "Re: Original",
						Description:  "Reply Content",
						Sending_date: testTime,
						ParentID:     1,
					}, nil)
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{
						ID:           1,
						Sender_email: "recipient@example.com",
						Recipient:    "test@example.com",
						Title:        "Original",
						Description:  "Original Content",
						Sending_date: testTime,
						ParentID:     0,
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: []models.Email{
				{
					ID:           2,
					Sender_email: "test@example.com",
					Recipient:    "recipient@example.com",
					Title:        "Re: Original",
					Description:  "Reply Content",
					Sending_date: testTime,
					ParentID:     1,
				},
				{
					ID:           1,
					Sender_email: "recipient@example.com",
					Recipient:    "test@example.com",
					Title:        "Original",
					Description:  "Original Content",
					Sending_date: testTime,
					ParentID:     0,
				},
			},
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
			name:      "письмо не найдено",
			setupAuth: true,
			emailID:   "999",
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(999).
					Return(models.Email{}, errors.New("not found"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "email_not_found",
			},
		},
		{
			name:      "доступ к чужому письму",
			setupAuth: true,
			emailID:   "1",
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{
						ID:           1,
						Sender_email: "other@example.com",
						Recipient:    "another@example.com",
						Title:        "Private Email",
					}, nil)
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: models.Error{
				Status: http.StatusUnauthorized,
				Body:   "unauthorized",
			},
		},
		{
			name:       "некорректный ID в пути",
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

			req := httptest.NewRequest(http.MethodGet, "/email/"+tt.emailID, nil)

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
			router.SingleEmailHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response interface{}
			if tt.wantStatus == http.StatusOK {
				var emails []models.Email
				err := json.NewDecoder(w.Body).Decode(&emails)
				assert.NoError(t, err)
				response = emails
			} else {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				response = errResponse
			}

			assert.Equal(t, tt.wantBody, response)
			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}
		})
	}
}
