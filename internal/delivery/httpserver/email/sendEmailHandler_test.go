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
		SaveEmail(gomock.Any()).
		Return(nil)

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

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(bodyBytes))
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
		SaveEmail(gomock.Any()).
		Return(nil)

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

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(bodyBytes))
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
		SaveEmail(gomock.Any()).
		Return(nil)

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

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(bodyBytes))
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

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(bodyBytes))
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

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer([]byte(`{"invalid json`)))
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

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	rr := httptest.NewRecorder()

	emailRouter.SendEmailHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSendEmailHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	req := converters.SendEmailRequest{
		ParentId:    0,
		Recipient:   "recipient@example.com",
		Title:       "Test Email",
		Description: "Test Description",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		SaveEmail(gomock.Any()).
		Return(nil)

	mockEmailUseCase.EXPECT().
		SendEmail("sender@example.com", []string{"recipient@example.com"}, req.Title, req.Description).
		Return(nil)

	router.SendEmailHandler(rr, httpReq)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSendEmailHandler_Reply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@example.com",
		Title:        "Original Email",
		Description:  "Original Description",
	}

	req := converters.SendEmailRequest{
		ParentId:    1,
		Title:       "Re: Original Email",
		Description: "Reply Description",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(1).
		Return(originalEmail, nil)

	mockEmailUseCase.EXPECT().
		SaveEmail(gomock.Any()).
		Return(nil)

	mockEmailUseCase.EXPECT().
		ReplyEmail("sender@example.com", "original@example.com", originalEmail, req.Description).
		Return(nil)

	router.SendEmailHandler(rr, httpReq)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSendEmailHandler_Forward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@example.com",
		Description:  "Original Description",
	}

	req := converters.SendEmailRequest{
		ParentId:    1,
		Recipient:   "forward@example.com",
		Title:       "Fwd: Original Email",
		Description: "Forward Description",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(1).
		Return(originalEmail, nil)

	mockEmailUseCase.EXPECT().
		SaveEmail(gomock.Any()).
		Return(nil)

	mockEmailUseCase.EXPECT().
		ForwardEmail("sender@example.com", []string{"forward@example.com"}, originalEmail).
		Return(nil)

	router.SendEmailHandler(rr, httpReq)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSendEmailHandler_InvalidOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@example.com",
	}

	req := converters.SendEmailRequest{
		ParentId:    1,
		Title:       "Invalid Operation",
		Description: "Description",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(1).
		Return(originalEmail, nil)

	router.SendEmailHandler(rr, httpReq)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response models.Error
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "invalid_operation", response.Body)
}

func TestSendEmailHandler_Errors(t *testing.T) {
	testCases := []struct {
		name           string
		setupRequest   func() (*http.Request, *gomock.Controller)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Unauthorized",
			setupRequest: func() (*http.Request, *gomock.Controller) {
				ctrl := gomock.NewController(t)
				req := httptest.NewRequest(http.MethodPost, "/email", nil)
				return req, ctrl
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "unauthorized",
		},
		{
			name: "Invalid Request Body",
			setupRequest: func() (*http.Request, *gomock.Controller) {
				ctrl := gomock.NewController(t)
				req := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBufferString("invalid json"))
				ctx := context.WithValue(req.Context(), "email", "sender@example.com")
				req = req.WithContext(ctx)
				return req, ctrl
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request_body",
		},
		{
			name: "Parent Email Not Found",
			setupRequest: func() (*http.Request, *gomock.Controller) {
				ctrl := gomock.NewController(t)
				req := converters.SendEmailRequest{
					ParentId: 999,
					Title:    "Re: Not Found",
				}
				reqBody, _ := json.Marshal(req)
				httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
				ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
				httpReq = httpReq.WithContext(ctx)
				return httpReq, ctrl
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "parent_email_not_found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpReq, ctrl := tc.setupRequest()
			defer ctrl.Finish()

			mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
			router := NewEmailRouter(mockEmailUseCase)

			if tc.name == "Parent Email Not Found" {
				mockEmailUseCase.EXPECT().
					GetEmailByID(999).
					Return(models.Email{}, assert.AnError)
			}

			rr := httptest.NewRecorder()
			router.SendEmailHandler(rr, httpReq)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			var response models.Error
			json.NewDecoder(rr.Body).Decode(&response)
			assert.Equal(t, tc.expectedError, response.Body)
		})
	}
}

func TestSendEmailHandler_SaveReplyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@example.com",
		Title:        "Original Email",
	}

	req := converters.SendEmailRequest{
		ParentId:    1,
		Title:       "Re: Original Email",
		Description: "Reply Description",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(1).
		Return(originalEmail, nil)

	mockEmailUseCase.EXPECT().
		SaveEmail(gomock.Any()).
		Return(assert.AnError)

	router.SendEmailHandler(rr, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var response models.Error
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "failed_to_save_reply", response.Body)
}

func TestSendEmailHandler_SaveForwardError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	originalEmail := models.Email{
		ID:           1,
		Sender_email: "original@example.com",
		Description:  "Original Description",
	}

	req := converters.SendEmailRequest{
		ParentId:    1,
		Recipient:   "forward@example.com",
		Title:       "Fwd: Original Email",
		Description: "Forward Description",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/email", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(httpReq.Context(), "email", "sender@example.com")
	httpReq = httpReq.WithContext(ctx)
	rr := httptest.NewRecorder()

	mockEmailUseCase.EXPECT().
		GetEmailByID(1).
		Return(originalEmail, nil)

	mockEmailUseCase.EXPECT().
		SaveEmail(gomock.Any()).
		Return(assert.AnError)

	router.SendEmailHandler(rr, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var response models.Error
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "failed_to_save_forward", response.Body)
}
