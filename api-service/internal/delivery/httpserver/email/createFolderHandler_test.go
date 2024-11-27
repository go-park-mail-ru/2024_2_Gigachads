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
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_CreateFolderHandler(t *testing.T) {
	tests := []struct {
		name       string
		input      models2.Folder
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   string
	}{
		{
			name: "успешное создание папки",
			input: models2.Folder{
				Name: "TestFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateFolder("test@example.com", "TestFolder").
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "неавторизованный запрос",
			input: models2.Folder{
				Name: "TestFolder",
			},
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name: "пустое имя папки",
			input: models2.Folder{
				Name: "",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateFolder(gomock.Any(), "").
					Return(errors.New("invalid folder name")).
					AnyTimes()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "error_with_creating_folder",
		},
		{
			name: "слишком длинное имя папки",
			input: models2.Folder{
				Name: strings.Repeat("a", 51),
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateFolder(gomock.Any(), strings.Repeat("a", 51)).
					Return(errors.New("invalid folder name")).
					AnyTimes()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "error_with_creating_folder",
		},
		{
			name: "ошибка создания папки",
			input: models2.Folder{
				Name: "TestFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateFolder("test@example.com", "TestFolder").
					Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "error_with_creating_folder",
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
			req := httptest.NewRequest(http.MethodPost, "/folders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.CreateFolderHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != "" {
				var response models2.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body)
			}
		})
	}
}
