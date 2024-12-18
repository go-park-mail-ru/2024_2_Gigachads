package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMultipartFormData(t *testing.T, fieldName, fileName string, fileContent []byte) (bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	part, err := w.CreateFormFile(fieldName, fileName)
	assert.NoError(t, err)
	_, err = io.Copy(part, bytes.NewReader(fileContent))
	assert.NoError(t, err)

	err = w.WriteField("name", fileName)
	assert.NoError(t, err)

	err = w.Close()
	assert.NoError(t, err)

	return b, w.FormDataContentType()
}

func TestEmailRouter_UploadAttachHandler(t *testing.T) {
	tests := []struct {
		name          string
		setupAuth     bool
		fileContent   []byte
		fileName      string
		mockSetup     func(*mocks.MockEmailUseCase)
		wantStatus    int
		wantBody      interface{}
		skipMultipart bool
		oversizedFile bool
	}{
		{
			name:        "успешная загрузка файла",
			setupAuth:   true,
			fileContent: []byte("test file content"),
			fileName:    "test.txt",
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UploadAttach(gomock.Any(), []byte("test file content"), "test.txt").
					Return("uploads/test.txt", nil)
			},
			wantStatus: http.StatusOK,
			wantBody: models.FilePath{
				Path: "uploads/test.txt",
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
			name:          "некорректный формат запроса",
			setupAuth:     true,
			skipMultipart: true,
			wantStatus:    http.StatusBadRequest,
			wantBody: models.Error{
				Status: http.StatusBadRequest,
				Body:   "error_with_parsing_file",
			},
		},
		{
			name:          "слишком большой файл",
			setupAuth:     true,
			oversizedFile: true,
			fileContent:   bytes.Repeat([]byte("a"), 11*1024*1024), // 11MB
			fileName:      "large.txt",
			wantStatus:    http.StatusRequestEntityTooLarge,
			wantBody: models.Error{
				Status: http.StatusRequestEntityTooLarge,
				Body:   "too_big_body",
			},
		},
		{
			name:        "ошибка при загрузке файла",
			setupAuth:   true,
			fileContent: []byte("test content"),
			fileName:    "error.txt",
			mockSetup: func(m *mocks.MockEmailUseCase) {
				m.EXPECT().
					UploadAttach(gomock.Any(), []byte("test content"), "error.txt").
					Return("", errors.New("upload error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: models.Error{
				Status: http.StatusInternalServerError,
				Body:   "error_with_upload_attach",
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

			var req *http.Request
			if tt.skipMultipart {
				req = httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader("invalid"))
			} else {
				body, contentType := createMultipartFormData(t, "file", tt.fileName, tt.fileContent)
				req = httptest.NewRequest(http.MethodPost, "/upload", &body)
				req.Header.Set("Content-Type", contentType)
			}

			if tt.setupAuth {
				ctx := context.WithValue(req.Context(), "email", "test@example.com")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.UploadAttachHandler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response interface{}
			if tt.wantStatus == http.StatusOK {
				var filePath models.FilePath
				err := json.NewDecoder(w.Body).Decode(&filePath)
				assert.NoError(t, err)
				response = filePath
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
