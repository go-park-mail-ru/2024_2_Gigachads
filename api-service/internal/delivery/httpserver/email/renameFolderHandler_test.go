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

func TestEmailRouter_RenameFolderHandler(t *testing.T) {
	tests := []struct {
		name       string
		setupAuth  bool
		folder     models.RenameFolder
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
		rawInput   string
	}{
		{
			name:      "успешное переименование папки",
			setupAuth: true,
			folder: models.RenameFolder{
				Name:    "OldFolder",
				NewName: "NewFolder",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					RenameFolder("test@example.com", "OldFolder", "NewFolder").
					Return(nil)
			},
			wantStatus: http.StatusOK,
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
			name:       "некорректный JSON в запросе",
			setupAuth:  true,
			rawInput:   "{invalid json",
			wantStatus: http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "invalid_json",
			},
		},
		{
			name:      "ошибка при переименовании папки",
			setupAuth: true,
			folder: models.RenameFolder{
				Name:    "OldFolder",
				NewName: "NewFolder",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					RenameFolder("test@example.com", "OldFolder", "NewFolder").
					Return(errors.New("rename error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "error_with_rename_folder",
			},
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
			if tt.rawInput != "" {
				reqBody = []byte(tt.rawInput)
			} else {
				reqBody, err = json.Marshal(tt.folder)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/renamefolder", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.RenameFolderHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus != http.StatusOK {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, errResponse)
			} else {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}
		})
	}
}
