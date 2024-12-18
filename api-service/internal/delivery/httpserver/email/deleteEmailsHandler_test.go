package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_DeleteEmailsHandler(t *testing.T) {
	tests := []struct {
		name        string
		input       DeleteEmailsRequest
		setupAuth   bool
		mockSetup   func(*mocks.MockEmailUseCase)
		wantStatus  int
		wantBody    string
		useRawInput bool
		rawInput    string
	}{
		{
			name: "успешное удаление писем",
			input: DeleteEmailsRequest{
				IDs:    []string{"1", "2", "3"},
				Folder: "Inbox",
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
			name:        "некорректный JSON",
			rawInput:    `{"ids": [1,2,3], "folder": }`,
			setupAuth:   true,
			useRawInput: true,
			wantStatus:  http.StatusBadRequest,
			wantBody:    "неверный формат данных",
		},
		{
			name: "пустой список ID",
			input: DeleteEmailsRequest{
				IDs:    []string{},
				Folder: "Inbox",
			},
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "список ID пуст",
		},
		{
			name: "некорректный формат ID",
			input: DeleteEmailsRequest{
				IDs:    []string{"1", "abc", "3"},
				Folder: "Inbox",
			},
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "неверный формат ID",
		},
		{
			name: "ошибка при удалении писем",
			input: DeleteEmailsRequest{
				IDs:    []string{"1", "2", "3"},
				Folder: "Inbox",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteEmails("test@example.com", []int{1, 2, 3}).
					Return(errors.New("db error"))
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

			var reqBody []byte
			var err error
			if tt.useRawInput {
				reqBody = []byte(tt.rawInput)
			} else {
				reqBody, err = json.Marshal(tt.input)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.DeleteEmailsHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != "" {
				var response struct {
					Body string `json:"body"`
				}
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body)
			}
		})
	}
}
