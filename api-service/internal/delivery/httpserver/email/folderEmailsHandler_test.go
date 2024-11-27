package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	models2 "mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_FolderEmailsHandler(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name        string
		input       models2.Folder
		setupAuth   bool
		mockSetup   func(*mocks.MockEmailUseCase)
		wantStatus  int
		wantBody    string
		wantData    []models2.Email
		useRawInput bool
		rawInput    string
	}{
		{
			name: "успешное получение писем",
			input: models2.Folder{
				Name: "Inbox",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolderEmails("test@example.com", "Inbox").
					Return([]models2.Email{
						{
							ID:           1,
							ParentID:     0,
							Sender_email: "sender@example.com",
							Recipient:    "test@example.com",
							Title:        "Test Email",
							Description:  "Test Content",
							IsRead:       false,
							Sending_date: testTime,
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantData: []models2.Email{
				{
					ID:           1,
					ParentID:     0,
					Sender_email: "sender@example.com",
					Recipient:    "test@example.com",
					Title:        "Test Email",
					Description:  "Test Content",
					IsRead:       false,
					Sending_date: testTime,
				},
			},
		},
		{
			name:       "неавторизованный запрос",
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:        "некорректный JSON",
			setupAuth:   true,
			useRawInput: true,
			rawInput:    `{"name": }`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    "invalid_json",
		},
		{
			name: "ошибка получения писем",
			input: models2.Folder{
				Name: "Inbox",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolderEmails("test@example.com", "Inbox").
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "email_not_found",
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
			if tt.useRawInput {
				reqBody = []byte(tt.rawInput)
			} else {
				reqBody, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(http.MethodPost, "/folder/emails", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.FolderEmailsHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")

			if tt.wantBody != "" {
				var response models2.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body, "Unexpected error body")
			}

			if tt.wantData != nil {
				var response []models2.Email
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)

				for i := range response {
					assert.Equal(t, tt.wantData[i].ID, response[i].ID)
					assert.Equal(t, tt.wantData[i].ParentID, response[i].ParentID)
					assert.Equal(t, tt.wantData[i].Sender_email, response[i].Sender_email)
					assert.Equal(t, tt.wantData[i].Recipient, response[i].Recipient)
					assert.Equal(t, tt.wantData[i].Title, response[i].Title)
					assert.Equal(t, tt.wantData[i].Description, response[i].Description)
					assert.Equal(t, tt.wantData[i].IsRead, response[i].IsRead)
				}
			}
		})
	}
}
