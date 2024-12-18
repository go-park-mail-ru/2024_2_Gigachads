package email

import (
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

func TestEmailRouter_FoldersHandler(t *testing.T) {
	tests := []struct {
		name       string
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
	}{
		{
			name:      "успешное получение папок",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolders("test@example.com").
					Return([]string{"Inbox", "Sent", "Custom"}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				"Inbox",
				"Sent",
				"Custom",
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
			name:      "ошибка при получении папок",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolders("test@example.com").
					Return(nil, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "error_with_getting_folders",
			},
		},
		{
			name:      "пустой список папок",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolders("test@example.com").
					Return([]string{}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []string{},
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

			req := httptest.NewRequest(http.MethodGet, "/folder", nil)
			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.FoldersHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response interface{}
			if tt.wantStatus == http.StatusOK {
				var folders []string
				err := json.NewDecoder(w.Body).Decode(&folders)
				assert.NoError(t, err)
				response = folders
			} else {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				response = errResponse
			}

			assert.Equal(t, tt.wantBody, response)
			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}
		})
	}
}
