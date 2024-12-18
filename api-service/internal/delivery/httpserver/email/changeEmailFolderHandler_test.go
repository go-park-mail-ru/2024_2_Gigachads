package email

import (
	"context"
	"encoding/json"
	"errors"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmailRouter_ChangeEmailFolderHandler(t *testing.T) {
	tests := []struct {
		name       string
		emailID    string
		folder     models.Folder
		setupAuth  bool
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   string
	}{
		{
			name:    "успешное изменение папки",
			emailID: "1",
			folder: models.Folder{
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
			name:       "неавторизованный запрос",
			emailID:    "1",
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:    "некорректный ID письма",
			emailID: "invalid",
			folder: models.Folder{
				Name: "NewFolder",
			},
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid_path",
		},
		{
			name:    "ошибка при изменении папки",
			emailID: "1",
			folder: models.Folder{
				Name: "NewFolder",
			},
			setupAuth: true,
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					ChangeEmailFolder(1, "test@example.com", "NewFolder").
					Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cant_change_name",
		},
		{
			name:       "некорректный JSON в запросе",
			emailID:    "1",
			setupAuth:  true,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid_json",
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

			var reqBody string
			if tt.folder.Name != "" {
				folderJSON, err := json.Marshal(tt.folder)
				assert.NoError(t, err)
				reqBody = string(folderJSON)
			} else {
				reqBody = "invalid json"
			}

			req := httptest.NewRequest(http.MethodPut, "/emails/"+tt.emailID+"/folder", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			vars := map[string]string{
				"id": tt.emailID,
			}
			req = mux.SetURLVars(req, vars)

			if tt.setupAuth {
				req = req.WithContext(NewContextWithEmail(req.Context(), "test@example.com"))
			}

			w := httptest.NewRecorder()
			router.ChangeEmailFolderHandler(w, req)

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

func NewContextWithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, "email", email)
}
