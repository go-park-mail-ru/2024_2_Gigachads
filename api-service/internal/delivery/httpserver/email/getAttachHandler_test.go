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

func TestEmailRouter_GetAttachHandler(t *testing.T) {
	tests := []struct {
		name       string
		setupAuth  bool
		filePath   models.FilePath
		mockSetup  func(*mocks.MockEmailUseCase)
		wantStatus int
		wantBody   interface{}
		rawInput   string
		wantData   []byte
	}{
		{
			name:      "успешное получение файла",
			setupAuth: true,
			filePath: models.FilePath{
				Path: "test/path/file.txt",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetAttach(gomock.Any(), "test/path/file.txt").
					Return([]byte("test file content"), nil)
			},
			wantStatus: http.StatusOK,
			wantData:   []byte("test file content"),
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
			name:      "ошибка при получении файла",
			setupAuth: true,
			filePath: models.FilePath{
				Path: "test/path/nonexistent.txt",
			},
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					GetAttach(gomock.Any(), "test/path/nonexistent.txt").
					Return(nil, errors.New("file not found"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "failed_to_get_file",
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
				reqBody, err = json.Marshal(tt.filePath)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/getattach", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.GetAttachHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "multipart/form-data", w.Header().Get("Content-Type"))
				assert.Equal(t, tt.wantData, w.Body.Bytes())
			} else {
				var errResponse models.Error
				err := json.NewDecoder(w.Body).Decode(&errResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, errResponse)
			}
		})
	}
}
