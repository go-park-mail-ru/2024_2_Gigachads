package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mail/internal/delivery/converters"
	"mail/internal/delivery/httpserver/email/mocks"
	"mail/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source=../../../models/email_model.go -destination=mocks/mock_email_usecase.go -package=mocks

func TestSendEmailHandler_SendSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockUseCase.EXPECT().
		SendEmail("sender@example.com", []string{"recipient@example.com"}, "Test Subject", "Test Body").
		Return(nil)

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		ParentId:    0,
		Recipient:   "recipient@example.com",
		Title:       "Test Subject",
		Description: "Test Body",
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

func TestSendEmailHandler_ForwardSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@example.com",
		Recipient:    "sender@example.com",
		Title:        "Original Subject",
		Description:  "Original Body",
		Sending_date: time.Now(),
	}

	mockUseCase.EXPECT().
		GetEmailByID(1).
		Return(originalEmail, nil)

	mockUseCase.EXPECT().
		ForwardEmail("sender@example.com", []string{"forward@example.com"}, originalEmail).
		Return(nil)

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		ParentId:    1,
		Recipient:   "forward@example.com",
		Title:       "Fwd: Original Subject",
		Description: "Forwarding the original email.",
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

func TestSendEmailHandler_ReplySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	originalEmail := models.Email{
		ID:           2,
		Sender_email: "original@example.com",
		Recipient:    "sender@example.com",
		Title:        "Original Subject",
		Description:  "Original Body",
		Sending_date: time.Now(),
	}

	mockUseCase.EXPECT().
		GetEmailByID(2).
		Return(originalEmail, nil)

	mockUseCase.EXPECT().
		ReplyEmail("sender@example.com", "original@example.com", originalEmail, "This is a reply.").
		Return(nil)

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		ParentId:    2,
		Recipient:   "original@example.com",
		Title:       "Re: Original Subject",
		Description: "This is a reply.",
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

func TestSendEmailHandler_InvalidParentId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockUseCase.EXPECT().
		GetEmailByID(999).
		Return(models.Email{}, errors.New("email not found"))

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		ParentId:    999,
		Recipient:   "forward@example.com",
		Title:       "Fwd: Non-existent Email",
		Description: "Trying to forward a non-existent email.",
	}
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/send-email", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "email", "sender@example.com"))
	rr := httptest.NewRecorder()

	emailRouter.SendEmailHandler(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSendEmailHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEmailUseCase(ctrl)
	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}

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
		ParentId:    0,
		Recipient:   "recipient@example.com",
		Title:       "Test Subject",
		Description: "Test Body",
	}
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

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
		SendEmail("sender@example.com", []string{"recipient@example.com"}, "Test Subject", "Test Body").
		Return(errors.New("smtp error"))

	emailRouter := &EmailRouter{EmailUseCase: mockUseCase}
	requestBody := converters.SendEmailRequest{
		ParentId:    0,
		Recipient:   "recipient@example.com",
		Title:       "Test Subject",
		Description: "Test Body",
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
