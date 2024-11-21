package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEmailsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	request := DeleteEmailsRequest{
		IDs:    []string{"1", "2", "3"},
		Folder: "inbox",
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		DeleteEmails("test@example.com", []int{1, 2, 3}, "inbox").
		Return(nil)

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteEmailsHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodDelete, "/emails", nil)
	rr := httptest.NewRecorder()

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", response.Body)
}

func TestDeleteEmailsHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	invalidJSON := []byte(`{"ids": [1, 2,`)
	req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(invalidJSON))
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "неверный формат данных", response.Body)
}

func TestDeleteEmailsHandler_EmptyIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	request := DeleteEmailsRequest{
		IDs:    []string{},
		Folder: "inbox",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "список ID пуст", response.Body)
}

func TestDeleteEmailsHandler_EmptyFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	request := DeleteEmailsRequest{
		IDs:    []string{"1", "2"},
		Folder: "",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "не указана папка", response.Body)
}

func TestDeleteEmailsHandler_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	request := DeleteEmailsRequest{
		IDs:    []string{"1", "invalid", "3"},
		Folder: "inbox",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "неверный формат ID", response.Body)
}

func TestDeleteEmailsHandler_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	request := DeleteEmailsRequest{
		IDs:    []string{"1", "2", "3"},
		Folder: "inbox",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodDelete, "/emails", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		DeleteEmails("test@example.com", []int{1, 2, 3}, "inbox").
		Return(errors.New("delete error"))

	router.DeleteEmailsHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "ошибка при удалении писем", response.Body)
}
