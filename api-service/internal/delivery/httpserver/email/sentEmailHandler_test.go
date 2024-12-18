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
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_SentEmailsHandler(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
	}{
		{
			name:      "успешное получение отправленных писем",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetSentEmails("test@example.com").
					Return([]models.Email{
						{
							ID:           1,
							Sender_email: "test@example.com",
							Recipient:    "recipient@example.com",
							Title:        "Test Subject",
							Description:  "Test Body",
							Sending_date: testTime,
							IsRead:       true,
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: []models.Email{
				{
					ID:           1,
					Sender_email: "test@example.com",
					Recipient:    "recipient@example.com",
					Title:        "Test Subject",
					Description:  "Test Body",
					Sending_date: testTime,
					IsRead:       true,
				},
			},
		},
		{
			name:       "неавторизованный запрос",
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody: models.Error{
				Status: http.StatusUnauthorized,
				Body:   "unauthorized",
			},
		},
		{
			name:      "ошибка при получении писем",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetSentEmails("test@example.com").
					Return(nil, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "failed_to_get_sent_emails",
			},
		},
		{
			name:      "пустой список отправленных писем",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetSentEmails("test@example.com").
					Return([]models.Email{}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []models.Email{},
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

			req := httptest.NewRequest(http.MethodGet, "/sent", nil)

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.SentEmailsHandler(w, req)

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
