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

func TestEmailRouter_UpdateDraftHandler(t *testing.T) {
	tests := []struct {
		name        string
		setupAuth   bool
		input       models2.Draft
		mockSetup   func(*mocks.MockEmailUseCase)
		wantStatus  int
		wantBody    string
		useRawInput bool
		rawInput    string
	}{
		{
			name:      "успешное обновление черновика",
			setupAuth: true,
			input: models2.Draft{
				ID:          1,
				Recipient:   "recipient@example.com",
				Title:       "Updated Draft",
				Description: "Updated content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UpdateDraft(gomock.Any()).
					DoAndReturn(func(draft models2.Draft) error {
						assert.Equal(t, 1, draft.ID)
						assert.Equal(t, "recipient@example.com", draft.Recipient)
						assert.Equal(t, "Updated Draft", draft.Title)
						assert.Equal(t, "Updated content", draft.Description)
						return nil
					})
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
			setupAuth:   true,
			useRawInput: true,
			rawInput:    `{"id": 1, "title": }`,
			wantStatus:  http.StatusBadRequest,
			wantBody:    "invalid_json",
		},
		{
			name:      "ошибка при обновлении черновика",
			setupAuth: true,
			input: models2.Draft{
				ID:          1,
				Recipient:   "recipient@example.com",
				Title:       "Test Draft",
				Description: "Test content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UpdateDraft(gomock.Any()).
					Return(errors.New("update error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cant_update_draft",
		},
		{
			name:      "проверка санитизации данных",
			setupAuth: true,
			input: models2.Draft{
				ID:          1,
				Recipient:   "<script>alert('xss')</script>recipient@example.com",
				Title:       "<script>alert('xss')</script>Test Title",
				Description: "<script>alert('xss')</script>Test Content",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UpdateDraft(gomock.Any()).
					DoAndReturn(func(draft models2.Draft) error {
						assert.Equal(t, "recipient@example.com", draft.Recipient)
						assert.Equal(t, "Test Title", draft.Title)
						assert.Equal(t, "Test Content", draft.Description)
						return nil
					})
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
			if tt.useRawInput {
				reqBody = []byte(tt.rawInput)
			} else {
				reqBody, err = json.Marshal(tt.input)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/draft", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.UpdateDraftHandler(w, req)

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
