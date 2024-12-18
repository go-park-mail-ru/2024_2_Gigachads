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

func TestEmailRouter_SendDraftHandler(t *testing.T) {

	tests := []struct {
		name       string
		setupAuth  bool
		draft      models.Email
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
		rawInput   string
	}{
		{
			name:      "успешная отправка черновика",
			setupAuth: true,
			draft: models.Email{
				Recipient:   "recipient@example.com",
				Title:       "Draft Email",
				Description: "Draft Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(nil)
				m.EXPECT().
					SendEmail(
						gomock.Any(),
						"test@example.com",
						[]string{"recipient@example.com"},
						"Draft Email",
						"Draft Content",
					).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"status": "success",
			},
		},
		{
			name:      "успешный ответ из черновика",
			setupAuth: true,
			draft: models.Email{
				ParentID:    1,
				Title:       "Re: Original Email",
				Description: "Reply Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				originalEmail := models.Email{
					ID:           1,
					Sender_email: "original@example.com",
					Title:        "Original Email",
				}
				m.EXPECT().
					GetEmailByID(1).
					Return(originalEmail, nil)
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(nil)
				m.EXPECT().
					ReplyEmail(
						gomock.Any(),
						"test@example.com",
						"original@example.com",
						originalEmail,
						"Reply Content",
					).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"status": "success",
			},
		},
		{
			name:      "успешная пересылка из черновика",
			setupAuth: true,
			draft: models.Email{
				ParentID:    1,
				Title:       "Fwd: Original Email",
				Recipient:   "new@example.com,another@example.com",
				Description: "Forwarded Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				originalEmail := models.Email{
					ID:           1,
					Sender_email: "original@example.com",
					Title:        "Original Email",
				}
				m.EXPECT().
					GetEmailByID(1).
					Return(originalEmail, nil)
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(nil)
				m.EXPECT().
					ForwardEmail(
						gomock.Any(),
						"test@example.com",
						[]string{"new@example.com", "another@example.com"},
						originalEmail,
					).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"status": "success",
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
				Body:   "invalid_request_body",
			},
		},
		{
			name:      "ошибка при сохранении черновика",
			setupAuth: true,
			draft: models.Email{
				Recipient:   "recipient@example.com",
				Title:       "Draft Email",
				Description: "Draft Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(errors.New("save error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "cant_send_draft",
			},
		},
		{
			name:      "родительское письмо не найдено",
			setupAuth: true,
			draft: models.Email{
				ParentID:    999,
				Title:       "Re: Non-existent Email",
				Description: "Reply Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(999).
					Return(models.Email{}, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "parent_email_not_found",
			},
		},
		{
			name:      "ошибка при сохранении ответа",
			setupAuth: true,
			draft: models.Email{
				ParentID:    1,
				Title:       "Re: Original Email",
				Description: "Reply Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				originalEmail := models.Email{
					ID:           1,
					Sender_email: "original@example.com",
					Title:        "Original Email",
				}
				m.EXPECT().
					GetEmailByID(1).
					Return(originalEmail, nil)
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(errors.New("save error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "failed_to_save_reply",
			},
		},
		{
			name:      "ошибка при сохранении пересылки",
			setupAuth: true,
			draft: models.Email{
				ParentID:    1,
				Title:       "Fwd: Original Email",
				Description: "Forward Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				originalEmail := models.Email{
					ID:           1,
					Sender_email: "original@example.com",
					Title:        "Original Email",
				}
				m.EXPECT().
					GetEmailByID(1).
					Return(originalEmail, nil)
				m.EXPECT().
					SendDraft(gomock.Any()).
					Return(errors.New("save error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "failed_to_save_forward",
			},
		},
		{
			name:      "некорректная операция с родительским письмом",
			setupAuth: true,
			draft: models.Email{
				ParentID:    1,
				Title:       "Invalid Operation",
				Description: "Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetEmailByID(1).
					Return(models.Email{}, nil)
			},
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_operation",
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
				reqBody, err = json.Marshal(tt.draft)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/senddraft", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.SendDraftHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response interface{}
			if tt.wantStatus == http.StatusOK {
				var successResponse map[string]string
				err := json.NewDecoder(w.Body).Decode(&successResponse)
				assert.NoError(t, err)
				response = successResponse
			} else {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				response = errResponse
			}

			assert.Equal(t, tt.wantBody, response)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}
