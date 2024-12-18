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
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_CreateDraftHandler(t *testing.T) {
	tests := []struct {
		name        string
		input       models.Email
		setupAuth   bool
		mockSetup   func(*mocks.MockEmailUseCase)
		wantStatus  int
		wantBody    string
		useRawInput bool
		rawInput    string
	}{
		{
			name: "успешное создание черновика",
			input: models.Email{
				Recipient:   "recipient@example.com",
				Title:       "Test Draft",
				Description: "Test Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateDraft(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "создание черновика ответа",
			input: models.Email{
				ParentID:    1,
				Title:       "Re: Original Email",
				Description: "Reply Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{
						Sender_email: "original@example.com",
					}, nil)
				m.EXPECT().
					CreateDraft(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "создание черновика пересылки",
			input: models.Email{
				ParentID:    1,
				Recipient:   "forward@example.com",
				Title:       "Fwd: Original Email",
				Description: "Forward Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{}, nil)
				m.EXPECT().
					CreateDraft(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "неавторизованный запрос",
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:        "некорректный JSON",
			rawInput:    `{"recipient": "test@example.com", "title": }`,
			setupAuth:   true,
			useRawInput: true,
			wantStatus:  http.StatusBadRequest,
			wantBody:    "invalid_request_body",
		},
		{
			name: "ошибка получения родительского письма",
			input: models.Email{
				ParentID:    1,
				Title:       "Re: Original Email",
				Description: "Reply Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{}, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "parent_email_not_found",
		},
		{
			name: "ошибка создания черновика",
			input: models.Email{
				Recipient:   "recipient@example.com",
				Title:       "Test Draft",
				Description: "Test Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateDraft(gomock.Any()).
					Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "error_with_creating_draft",
		},
		{
			name: "некорректная операция с родительским письмом",
			input: models.Email{
				ParentID:    1,
				Title:       "Invalid Operation",
				Description: "Test Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{}, nil)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid_operation",
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
			if tt.useRawInput {
				reqBody = []byte(tt.rawInput)
			} else {
				reqBody, err = json.Marshal(tt.input)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/draft", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.CreateDraftHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != "" {
				var response models.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body)
			} else if w.Code == http.StatusOK {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
			}
		})
	}
}
