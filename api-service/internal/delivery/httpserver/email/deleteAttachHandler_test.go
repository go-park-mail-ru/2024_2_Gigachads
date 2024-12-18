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

func TestEmailRouter_DeleteAttachHandler(t *testing.T) {
	tests := []struct {
		name        string
		input       models.FilePath
		setupAuth   bool
		mockSetup   func(*mocks.MockEmailUseCase)
		wantStatus  int
		wantBody    string
		useRawInput bool
		rawInput    string
	}{
		{
			name: "успешное удаление вложения",
			input: models.FilePath{
				Path: "/path/to/file.pdf",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteAttach(gomock.Any(), "/path/to/file.pdf").
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
			rawInput:    `{"path": }`,
			setupAuth:   true,
			useRawInput: true,
			wantStatus:  http.StatusBadRequest,
			wantBody:    "invalid_json",
		},
		{
			name: "ошибка при удалении файла",
			input: models.FilePath{
				Path: "/path/to/file.pdf",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteAttach(gomock.Any(), "/path/to/file.pdf").
					Return(errors.New("failed to delete file"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed_to_delete_file",
		},
		{
			name: "пустой путь к файлу",
			input: models.FilePath{
				Path: "",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteAttach(gomock.Any(), "").
					Return(errors.New("invalid file path"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed_to_delete_file",
		},
		{
			name: "таймаут операции",
			input: models.FilePath{
				Path: "/path/to/file.pdf",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteAttach(gomock.Any(), "/path/to/file.pdf").
					DoAndReturn(func(ctx context.Context, path string) error {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case <-time.After(6 * time.Second):
							return nil
						}
					})
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed_to_delete_file",
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

			req := httptest.NewRequest(http.MethodDelete, "/attachments", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.DeleteAttachHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != "" {
				var response models.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body)
			}
		})
	}
}
