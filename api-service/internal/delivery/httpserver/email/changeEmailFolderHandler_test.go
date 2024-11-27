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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_ChangeEmailFolderHandler(t *testing.T) {
	tests := []struct {
		name       string
		emailID    string
		input      models2.Folder
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   string
	}{
		{
			name:    "успешное изменение папки",
			emailID: "1",
			input: models2.Folder{
				Name: "NewFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					ChangeEmailFolder(1, "test@example.com", "NewFolder").
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "неавторизованный запрос",
			emailID: "1",
			input: models2.Folder{
				Name: "NewFolder",
			},
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:    "некорректный ID папки",
			emailID: "invalid",
			input: models2.Folder{
				Name: "NewFolder",
			},
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid_path",
		},
		{
			name:    "ошибка изменения папки",
			emailID: "1",
			input: models2.Folder{
				Name: "NewFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					ChangeEmailFolder(1, "test@example.com", "NewFolder").
					Return(errors.New("error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cant_change_name",
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

			// Создаем новый роутер mux
			r := mux.NewRouter()
			r.HandleFunc("/emails/{id}/folder", router.ChangeEmailFolderHandler).Methods("PUT")

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/emails/"+tt.emailID+"/folder", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

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
