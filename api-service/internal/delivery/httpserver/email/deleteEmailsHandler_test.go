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

func TestEmailRouter_DeleteEmailsHandler(t *testing.T) {
	tests := []struct {
		name       string
		input      DeleteEmailsRequest
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   string
	}{
		{
			name: "успешное удаление",
			input: DeleteEmailsRequest{
				IDs: []string{"1", "2", "3"},
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteEmails("test@example.com", []int{1, 2, 3}).
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "неавторизованный запрос",
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name: "пустой список ID",
			input: DeleteEmailsRequest{
				IDs: []string{},
			},
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "список ID пуст",
		},
		{
			name: "некорректный ID",
			input: DeleteEmailsRequest{
				IDs: []string{"1", "invalid", "3"},
			},
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "неверный формат ID",
		},
		{
			name: "ошибка удаления",
			input: DeleteEmailsRequest{
				IDs: []string{"1", "2", "3"},
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteEmails("test@example.com", []int{1, 2, 3}).
					Return(errors.New("error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "ошибка при удалении писем",
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
			req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.DeleteEmailsHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")

			if tt.wantBody != "" {
				var response models.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body, "Unexpected response body")
			}
		})
	}
}
