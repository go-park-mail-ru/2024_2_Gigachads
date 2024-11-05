package user

import (
	"bytes"
	"context"
	"io"
	"mail/internal/delivery/httpserver/email/mocks"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMultipartFormData(t *testing.T, fieldName, fileName string, fileContent []byte) (*bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile(fieldName, fileName)
	assert.NoError(t, err)

	_, err = io.Copy(fw, bytes.NewReader(fileContent))
	assert.NoError(t, err)

	w.Close()

	return &b, w.FormDataContentType()
}

func TestUserRouter_UploadAvatarHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewUserRouter(mockUserUseCase)

	t.Run("успешная загрузка аватара", func(t *testing.T) {
		fileContent := []byte("fake-image-content")
		body, contentType := createMultipartFormData(t, "avatar", "test.jpg", fileContent)

		req := httptest.NewRequest("POST", "/upload-avatar", body)
		req.Header.Set("Content-Type", contentType)
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			ChangeAvatar(gomock.Any(), gomock.Any(), "test@example.com").
			Return(nil)

		router.UploadAvatarHandler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("неавторизованный запрос", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload-avatar", nil)
		w := httptest.NewRecorder()

		router.UploadAvatarHandler(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ошибка парсинга формы", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload-avatar", bytes.NewBuffer([]byte("invalid form data")))
		req.Header.Set("Content-Type", "multipart/form-data")
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		router.UploadAvatarHandler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("отсутствие файла в форме", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest("POST", "/upload-avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		router.UploadAvatarHandler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
