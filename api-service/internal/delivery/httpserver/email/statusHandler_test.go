package email

import (
	"bytes"
	"context"
	"encoding/json"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmailStatusHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	status := Status{Status: true}
	emailID := "1"

	statusJSON, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/email/"+emailID+"/status", bytes.NewBuffer(statusJSON))

	vars := map[string]string{
		"id": emailID,
	}
	req = mux.SetURLVars(req, vars)

	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().ChangeStatus(1, true).Return(nil)

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEmailStatusHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	status := Status{Status: true}
	emailID := "1"

	statusJSON, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/email/"+emailID+"/status", bytes.NewBuffer(statusJSON))

	vars := map[string]string{
		"id": emailID,
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", response.Body)
}

func TestEmailStatusHandler_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	status := Status{Status: true}
	emailID := "invalid"

	statusJSON, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/email/"+emailID+"/status", bytes.NewBuffer(statusJSON))

	vars := map[string]string{
		"id": emailID,
	}
	req = mux.SetURLVars(req, vars)

	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_path", response.Body)
}

func TestEmailStatusHandler_ErrorChangingStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	status := Status{Status: true}
	emailID := "1"

	statusJSON, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, "/email/"+emailID+"/status", bytes.NewBuffer(statusJSON))

	vars := map[string]string{
		"id": emailID,
	}
	req = mux.SetURLVars(req, vars)

	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().ChangeStatus(1, true).Return(assert.AnError)

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_status", response.Body)
}

func TestEmailStatusHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	invalidJSON := []byte(`{"status": invalid}`)
	req := httptest.NewRequest(http.MethodPut, "/email/1/status", bytes.NewBuffer(invalidJSON))

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_json", response.Body)
}

func TestEmailStatusHandler_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	status := Status{Status: true}
	statusJSON, _ := json.Marshal(status)

	req := httptest.NewRequest(http.MethodPut, "/email/status", bytes.NewBuffer(statusJSON))
	req = mux.SetURLVars(req, map[string]string{})

	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_path", response.Body)
}

func TestEmailStatusHandler_EmptyBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := httptest.NewRequest(http.MethodPut, "/email/1/status", nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	ctx := context.WithValue(req.Context(), "email", "test@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	router.EmailStatusHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.Error
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_json", response.Body)
}

func TestEmailStatusHandler_DifferentStatuses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	testCases := []struct {
		name      string
		status    bool
		expectErr bool
		mockSetup func()
	}{
		{
			name:      "Valid Read Status",
			status:    true,
			expectErr: false,
			mockSetup: func() {
				mockEmailUseCase.EXPECT().ChangeStatus(1, true).Return(nil)
			},
		},
		{
			name:      "Valid Unread Status",
			status:    false,
			expectErr: false,
			mockSetup: func() {
				mockEmailUseCase.EXPECT().ChangeStatus(1, false).Return(nil)
			},
		},
		{
			name:      "Invalid Status",
			status:    false,
			expectErr: true,
			mockSetup: func() {
				mockEmailUseCase.EXPECT().ChangeStatus(1, false).Return(assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := Status{Status: tc.status}
			statusJSON, _ := json.Marshal(status)
			req := httptest.NewRequest(http.MethodPut, "/email/1/status", bytes.NewBuffer(statusJSON))

			vars := map[string]string{
				"id": "1",
			}
			req = mux.SetURLVars(req, vars)

			ctx := context.WithValue(req.Context(), "email", "test@example.com")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			tc.mockSetup()

			router.EmailStatusHandler(rr, req)

			if tc.expectErr {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}
