package email

import (
	"context"
	"encoding/json"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSingleEmailHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	email := models.Email{
		ID:       2,
		ParentID: 0,
		Title:    "Test Email",
	}

	req := httptest.NewRequest(http.MethodGet, "/email/2", nil)
	vars := map[string]string{
		"id": "2",
	}
	req = mux.SetURLVars(req, vars)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(2).
		Return(email, nil)

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, email, response[0])
}

func TestSingleEmailHandler_WithParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	child := models.Email{
		ID:       2,
		ParentID: 1,
		Title:    "Child Email",
	}

	parent := models.Email{
		ID:       1,
		ParentID: 0,
		Title:    "Parent Email",
	}

	gomock.InOrder(
		mockEmailUseCase.EXPECT().
			GetEmailByID(2).
			Return(child, nil),
		mockEmailUseCase.EXPECT().
			GetEmailByID(1).
			Return(parent, nil),
	)

	req := httptest.NewRequest(http.MethodGet, "/email/2", nil)
	vars := map[string]string{
		"id": "2",
	}
	req = mux.SetURLVars(req, vars)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
	assert.Equal(t, child, response[0])
	assert.Equal(t, parent, response[1])
}

func TestSingleEmailHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/email/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", response.Body)
}

func TestSingleEmailHandler_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/email/invalid", nil)
	vars := map[string]string{
		"id": "invalid",
	}
	req = mux.SetURLVars(req, vars)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_path", response.Body)
}

func TestSingleEmailHandler_NoID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/email/", nil)
	req = mux.SetURLVars(req, map[string]string{})
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_path", response.Body)
}

func TestSingleEmailHandler_SingleEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	email := models.Email{
		ID:       1,
		ParentID: 0,
		Title:    "Single Email",
	}

	req := httptest.NewRequest(http.MethodGet, "/email/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().GetEmailByID(1).Return(email, nil)

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, email, response[0])
}

func TestSingleEmailHandler_EmailNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodGet, "/email/999", nil)
	vars := map[string]string{
		"id": "999",
	}
	req = mux.SetURLVars(req, vars)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(999).
		Return(models.Email{}, assert.AnError)

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "email_not_found", response.Body)
}

func TestSingleEmailHandler_ParentNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	child := models.Email{
		ID:       2,
		ParentID: 1,
		Title:    "Child Email",
	}

	req := httptest.NewRequest(http.MethodGet, "/email/2", nil)
	vars := map[string]string{
		"id": "2",
	}
	req = mux.SetURLVars(req, vars)
	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	gomock.InOrder(
		mockEmailUseCase.EXPECT().
			GetEmailByID(2).
			Return(child, nil),
		mockEmailUseCase.EXPECT().
			GetEmailByID(1).
			Return(models.Email{}, assert.AnError),
	)

	router.SingleEmailHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Email
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, child, response[0])
}
