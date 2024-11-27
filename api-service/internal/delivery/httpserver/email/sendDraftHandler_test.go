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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_SendDraftHandler(t *testing.T) {
	tests := []struct {
		name       string
		input      models2.Draft
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   string
	}{
		{
			name: "успешная отправка нового письма",
			input: models2.Draft{
				ID:          1,
				Recipient:   "recipient@example.com",
				Title:       "Test Email",
				Description: "Test Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(nil)
				m.EXPECT().
					SendEmail(gomock.Any(), "test@example.com", []string{"recipient@example.com"}, "Test Email", "Test Content")
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "успешная отправка ответа",
			input: models2.Draft{
				ID:          2,
				ParentID:    1,
				Title:       "Re: Original Email",
				Description: "Reply Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models2.Email{
						Sender_email: "original@example.com",
					}, nil)
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(nil)
				m.EXPECT().
					ReplyEmail(gomock.Any(), "test@example.com", "original@example.com", gomock.Any(), "Reply Content")
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "успешная пересылка",
			input: models2.Draft{
				ID:          3,
				ParentID:    1,
				Recipient:   "forward@example.com",
				Title:       "Fwd: Original Email",
				Description: "Forward Content",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models2.Email{}, nil)
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(nil)
				m.EXPECT().
					ForwardEmail(gomock.Any(), "test@example.com", []string{"forward@example.com"}, gomock.Any())
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
			name: "некорректный ParentID",
			input: models2.Draft{
				ParentID: 1,
				Title:    "Re: Original Email",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models2.Email{}, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "parent_email_not_found",
		},
		{
			name: "ошибка отправки черновика",
			input: models2.Draft{
				Recipient: "recipient@example.com",
				Title:     "Test Email",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(errors.New("error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cant_send_draft",
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

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/drafts/send", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.SendDraftHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")

			if tt.wantBody != "" {
				var response models2.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body, "Unexpected error body")
			}
		})
	}
}
