package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	models2 "mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/api-service/internal/delivery/httpserver/email/mocks"
)

func TestEmailRouter_DeleteFolderHandler(t *testing.T) {
	tests := []struct {
		name        string
		input       models2.Folder
		setupAuth   bool
		mockSetup   func(*mocks.MockEmailUseCase)
		wantStatus  int
		wantBody    string
		useRawInput bool
		rawInput    string
	}{
		{
			name: "успешное удаление папки",
			input: models2.Folder{
				Name: "TestFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteFolder("test@example.com", "TestFolder").
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
			name: "ошибка удаления папки",
			input: models2.Folder{
				Name: "TestFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteFolder("test@example.com", "TestFolder").
					Return(errors.New("error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "error_with_deleting_folder",
		},
		{
			name:        "некорректный JSON",
			rawInput:    `{"name": }`,
			setupAuth:   true,
			useRawInput: true,
			wantStatus:  http.StatusBadRequest,
			wantBody:    "invalid_json",
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

			req := httptest.NewRequest(http.MethodDelete, "/folders", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.DeleteFolderHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")

			if tt.wantBody != "" {
				var response models2.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body, "Unexpected response body")
			}
		})
	}
}
