package user

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"mail/api-service/internal/delivery/httpserver/email/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter_UploadAvatarHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	router := NewUserRouter(mockUserUseCase)

	t.Run("успешная загрузка аватара", func(t *testing.T) {
		// Создаем тестовый файл
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("fake-image-content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/upload-avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Добавляем email в контекст
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		// Устанавливаем ожидание
		mockUserUseCase.EXPECT().
			ChangeAvatar(gomock.Any(), "test@example.com").
			Return(nil)

		router.UploadAvatarHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("неавторизованный запрос", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload-avatar", nil)
		w := httptest.NewRecorder()

		router.UploadAvatarHandler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "unauthorized", response["body"])
	})

	t.Run("ошибка при загрузке файла", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload-avatar", nil)
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		router.UploadAvatarHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "error_with_parsing_file", response["body"])
	})

	t.Run("ошибка при сохранении аватара", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("fake-image-content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/upload-avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := context.WithValue(req.Context(), "email", "test@example.com")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserUseCase.EXPECT().
			ChangeAvatar(gomock.Any(), "test@example.com").
			Return(assert.AnError)

		router.UploadAvatarHandler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "error_with_downloading_avatar", response["body"])
	})
}
