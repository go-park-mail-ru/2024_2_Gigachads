package email

import (
	"bytes"
	"context"
	"encoding/json"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func addEmailToContext(r *http.Request, email string) *http.Request {
	ctx := context.WithValue(r.Context(), "email", email)
	return r.WithContext(ctx)
}

const testEmail = "test@example.com"

func TestEmailRouter_FoldersWithUser(t *testing.T) {
	tests := []struct {
		name       string
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		endpoint   string
		method     string
		body       interface{}
		wantStatus int
		wantBody   string
		wantData   interface{}
	}{
		{
			name:      "получение списка папок авторизованного пользователя",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetFolders(testEmail).
					Return([]string{"Inbox", "Sent"}, nil)
			},
			endpoint:   "/folder",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantData:   []string{"Inbox", "Sent"},
		},
		{
			name:       "неавторизованный запрос",
			setupAuth:  false,
			endpoint:   "/folder",
			method:     http.MethodGet,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:      "создание новой папки",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					CreateFolder(testEmail, "NewFolder").
					Return(nil)
			},
			endpoint:   "/folder",
			method:     http.MethodPost,
			body:       models.Folder{Name: "NewFolder"},
			wantStatus: http.StatusOK,
		},
		{
			name:      "удаление папки",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					DeleteFolder(testEmail, "OldFolder").
					Return(nil)
			},
			endpoint:   "/folder",
			method:     http.MethodDelete,
			body:       models.Folder{Name: "OldFolder"},
			wantStatus: http.StatusOK,
		},
		{
			name:      "переименование папки",
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					RenameFolder(testEmail, "OldFolder", "NewFolder").
					Return(nil)
			},
			endpoint:   "/folder",
			method:     http.MethodPut,
			body:       models.RenameFolder{Name: "OldFolder", NewName: "NewFolder"},
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

			router := mux.NewRouter()
			emailRouter := NewEmailRouter(mockEmailUseCase)
			emailRouter.ConfigureEmailRouter(router)

			var reqBody []byte
			var err error
			if tt.body != nil {
				reqBody, err = json.Marshal(tt.body)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer(reqBody))
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			if tt.setupAuth {
				req = addEmailToContext(req, testEmail)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")

			if tt.wantBody != "" {
				var response models.Error
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, response.Body, "Unexpected error body")
			}

			if tt.wantData != nil {
				switch v := tt.wantData.(type) {
				case []string:
					var folders []string
					err := json.NewDecoder(w.Body).Decode(&folders)
					assert.NoError(t, err)
					assert.ElementsMatch(t, v, folders, "Unexpected folders")
				}
			}
		})
	}
}
