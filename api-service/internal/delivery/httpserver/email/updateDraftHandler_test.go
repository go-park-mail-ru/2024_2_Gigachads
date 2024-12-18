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

func TestEmailRouter_UpdateDraftHandler(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

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
			name:      "успешное обновление черновика",
			setupAuth: true,
			draft: models.Email{
				ID:           1,
				Sender_email: "test@example.com",
				Recipient:    "recipient@example.com",
				Title:        "Draft Title",
				Description:  "Draft Content",
				Sending_date: testTime,
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UpdateDraft(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
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
			name:      "ошибка при обновлении черновика",
			setupAuth: true,
			draft: models.Email{
				ID:           1,
				Sender_email: "test@example.com",
				Recipient:    "recipient@example.com",
				Title:        "Draft Title",
				Description:  "Draft Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UpdateDraft(gomock.Any()).
					Return(errors.New("update error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "cant_update_draft",
			},
		},
		{
			name:      "проверка санитизации данных",
			setupAuth: true,
			draft: models.Email{
				ID:           1,
				Sender_email: "test@example.com",
				Recipient:    "<script>alert('xss')</script>recipient@example.com",
				Title:        "<script>alert('xss')</script>Title",
				Description:  "<script>alert('xss')</script>Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UpdateDraft(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
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

			req := httptest.NewRequest(http.MethodPut, "/updatedraft", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.UpdateDraftHandler(w, req)

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
