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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_InboxStatusHandler(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	lastModified := time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		setupAuth  bool
		timestamp  models.Timestamp
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
		rawInput   string
	}{
		{
			name:      "успешное получение обновлений",
			setupAuth: true,
			timestamp: models.Timestamp{
				LastModified: lastModified,
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					InboxStatus(gomock.Any(), "test@example.com", lastModified).
					Return([]models.Email{
						{
							ID:           1,
							Sender_email: "sender@example.com",
							Recipient:    "test@example.com",
							Title:        "New Email",
							Description:  "New Content",
							Sending_date: testTime,
							IsRead:       false,
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: []models.Email{
				{
					ID:           1,
					Sender_email: "sender@example.com",
					Recipient:    "test@example.com",
					Title:        "New Email",
					Description:  "New Content",
					Sending_date: testTime,
					IsRead:       false,
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
			name:       "некорректный JSON в запросе",
			setupAuth:  true,
			rawInput:   "{invalid json",
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_json",
			},
		},
		{
			name:      "нет изменений",
			setupAuth: true,
			timestamp: models.Timestamp{
				LastModified: lastModified,
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					InboxStatus(gomock.Any(), "test@example.com", lastModified).
					Return(nil, errors.New("not modified"))
			},
			wantStatus: http.StatusNotModified,
			wantBody: models.Error{
				Status: http.StatusNotModified,
				Body:   "not_modified",
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
				reqBody, err = json.Marshal(tt.timestamp)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/inbox/status", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.InboxStatusHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var emails []models.Email
				err := json.NewDecoder(w.Body).Decode(&emails)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, emails)
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			} else {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, errResponse)
			}
		})
	}
}
