package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mail/internal/delivery/converters"
	"mail/internal/delivery/httpserver/email/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source=../../../models/email_model.go -destination=mocks/mock_email_usecase.go -package=mocks

func TestSendEmailHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockUseCase.EXPECT().
		SendEmail("sender@example.com", []string{"recipient@example.com"}, "Test Subject", "Test Body").
		Return(nil)

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		To:    "recipient@example.com",
		Title: "Test Subject",
		Body:  "Test Body",
	}
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/send-email", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "sender@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.SendEmailHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSendEmailHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

	// Отправляем невалидный JSON
	req, err := http.NewRequest("POST", "/send-email", bytes.NewBuffer([]byte(`{"invalid json`)))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "sender@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.SendEmailHandler(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSendEmailHandler_MissingEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

	requestBody := converters.SendEmailRequest{
		To:    "recipient@example.com",
		Title: "Test Subject",
		Body:  "Test Body",
	}
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// Не добавляем email в контекст
	req, err := http.NewRequest("POST", "/send-email", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	rr := httptest.NewRecorder()

	emailRouter.SendEmailHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSendEmailHandler_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockUseCase.EXPECT().
		SendEmail(
			"sender@example.com",
			[]string{"recipient@example.com"},
			"Test Subject",
			"Test Body",
		).Return(errors.New("smtp error"))

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		To:    "recipient@example.com",
		Title: "Test Subject",
		Body:  "Test Body",
	}
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/send-email", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "sender@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.SendEmailHandler(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
