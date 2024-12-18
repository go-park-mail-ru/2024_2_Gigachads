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

func TestEmailRouter_FolderEmailsHandler(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		setupAuth  bool
		folder     models.Folder
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
		rawInput   string
	}{
		{
			name:      "успешное получение писем из папки",
			setupAuth: true,
			folder: models.Folder{
				Name: "Inbox",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolderEmails("test@example.com", "Inbox").
					Return([]models.Email{
						{
							ID:           1,
							Sender_email: "sender@example.com",
							Recipient:    "test@example.com",
							Title:        "Test Subject",
							Description:  "Test Body",
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
					Title:        "Test Subject",
					Description:  "Test Body",
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
			name:      "ошибка при получении писем",
			setupAuth: true,
			folder: models.Folder{
				Name: "Custom",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolderEmails("test@example.com", "Custom").
					Return(nil, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "email_not_found",
			},
		},
		{
			name:      "пустой список писем",
			setupAuth: true,
			folder: models.Folder{
				Name: "Empty",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolderEmails("test@example.com", "Empty").
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

			var reqBody []byte
			var err error
			if tt.rawInput != "" {
				reqBody = []byte(tt.rawInput)
			} else if tt.folder.Name != "" {
				reqBody, err = json.Marshal(tt.folder)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/getfolder", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.FolderEmailsHandler(w, req)

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
