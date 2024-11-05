// Code generated by MockGen. DO NOT EDIT.
// Source: internal/models/email_model.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "mail/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailUseCase is a mock of EmailUseCase interface.
type MockEmailUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockEmailUseCaseMockRecorder
}

// MockEmailUseCaseMockRecorder is the mock recorder for MockEmailUseCase.
type MockEmailUseCaseMockRecorder struct {
	mock *MockEmailUseCase
}

// NewMockEmailUseCase creates a new mock instance.
func NewMockEmailUseCase(ctrl *gomock.Controller) *MockEmailUseCase {
	mock := &MockEmailUseCase{ctrl: ctrl}
	mock.recorder = &MockEmailUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailUseCase) EXPECT() *MockEmailUseCaseMockRecorder {
	return m.recorder
}

// ChangeStatus mocks base method.
func (m *MockEmailUseCase) ChangeStatus(id int, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeStatus", id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeStatus indicates an expected call of ChangeStatus.
func (mr *MockEmailUseCaseMockRecorder) ChangeStatus(id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeStatus", reflect.TypeOf((*MockEmailUseCase)(nil).ChangeStatus), id, status)
}

// DeleteEmails mocks base method.
func (m *MockEmailUseCase) DeleteEmails(userEmail string, messageIDs []int, folder string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmails", userEmail, messageIDs, folder)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmails indicates an expected call of DeleteEmails.
func (mr *MockEmailUseCaseMockRecorder) DeleteEmails(userEmail, messageIDs, folder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmails", reflect.TypeOf((*MockEmailUseCase)(nil).DeleteEmails), userEmail, messageIDs, folder)
}

// FetchEmailsViaPOP3 mocks base method.
func (m *MockEmailUseCase) FetchEmailsViaPOP3() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchEmailsViaPOP3")
	ret0, _ := ret[0].(error)
	return ret0
}

// FetchEmailsViaPOP3 indicates an expected call of FetchEmailsViaPOP3.
func (mr *MockEmailUseCaseMockRecorder) FetchEmailsViaPOP3() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchEmailsViaPOP3", reflect.TypeOf((*MockEmailUseCase)(nil).FetchEmailsViaPOP3))
}

// ForwardEmail mocks base method.
func (m *MockEmailUseCase) ForwardEmail(from string, to []string, originalEmail models.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForwardEmail", from, to, originalEmail)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForwardEmail indicates an expected call of ForwardEmail.
func (mr *MockEmailUseCaseMockRecorder) ForwardEmail(from, to, originalEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardEmail", reflect.TypeOf((*MockEmailUseCase)(nil).ForwardEmail), from, to, originalEmail)
}

// GetEmailByID mocks base method.
func (m *MockEmailUseCase) GetEmailByID(id int) (models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailByID", id)
	ret0, _ := ret[0].(models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmailByID indicates an expected call of GetEmailByID.
func (mr *MockEmailUseCaseMockRecorder) GetEmailByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailByID", reflect.TypeOf((*MockEmailUseCase)(nil).GetEmailByID), id)
}

// GetSentEmails mocks base method.
func (m *MockEmailUseCase) GetSentEmails(senderEmail string) ([]models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSentEmails", senderEmail)
	ret0, _ := ret[0].([]models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSentEmails indicates an expected call of GetSentEmails.
func (mr *MockEmailUseCaseMockRecorder) GetSentEmails(senderEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSentEmails", reflect.TypeOf((*MockEmailUseCase)(nil).GetSentEmails), senderEmail)
}

// Inbox mocks base method.
func (m *MockEmailUseCase) Inbox(id string) ([]models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inbox", id)
	ret0, _ := ret[0].([]models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inbox indicates an expected call of Inbox.
func (mr *MockEmailUseCaseMockRecorder) Inbox(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inbox", reflect.TypeOf((*MockEmailUseCase)(nil).Inbox), id)
}

// ReplyEmail mocks base method.
func (m *MockEmailUseCase) ReplyEmail(from, to string, originalEmail models.Email, replyText string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplyEmail", from, to, originalEmail, replyText)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplyEmail indicates an expected call of ReplyEmail.
func (mr *MockEmailUseCaseMockRecorder) ReplyEmail(from, to, originalEmail, replyText interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplyEmail", reflect.TypeOf((*MockEmailUseCase)(nil).ReplyEmail), from, to, originalEmail, replyText)
}

// SaveEmail mocks base method.
func (m *MockEmailUseCase) SaveEmail(email models.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEmail", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEmail indicates an expected call of SaveEmail.
func (mr *MockEmailUseCaseMockRecorder) SaveEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEmail", reflect.TypeOf((*MockEmailUseCase)(nil).SaveEmail), email)
}

// SendEmail mocks base method.
func (m *MockEmailUseCase) SendEmail(from string, to []string, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", from, to, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailUseCaseMockRecorder) SendEmail(from, to, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailUseCase)(nil).SendEmail), from, to, subject, body)
}

// MockEmailRepository is a mock of EmailRepository interface.
type MockEmailRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmailRepositoryMockRecorder
}

// MockEmailRepositoryMockRecorder is the mock recorder for MockEmailRepository.
type MockEmailRepositoryMockRecorder struct {
	mock *MockEmailRepository
}

// NewMockEmailRepository creates a new mock instance.
func NewMockEmailRepository(ctrl *gomock.Controller) *MockEmailRepository {
	mock := &MockEmailRepository{ctrl: ctrl}
	mock.recorder = &MockEmailRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailRepository) EXPECT() *MockEmailRepositoryMockRecorder {
	return m.recorder
}

// ChangeStatus mocks base method.
func (m *MockEmailRepository) ChangeStatus(id int, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeStatus", id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeStatus indicates an expected call of ChangeStatus.
func (mr *MockEmailRepositoryMockRecorder) ChangeStatus(id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeStatus", reflect.TypeOf((*MockEmailRepository)(nil).ChangeStatus), id, status)
}

// DeleteEmails mocks base method.
func (m *MockEmailRepository) DeleteEmails(userEmail string, messageIDs []int, folder string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmails", userEmail, messageIDs, folder)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmails indicates an expected call of DeleteEmails.
func (mr *MockEmailRepositoryMockRecorder) DeleteEmails(userEmail, messageIDs, folder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmails", reflect.TypeOf((*MockEmailRepository)(nil).DeleteEmails), userEmail, messageIDs, folder)
}

// GetEmailByID mocks base method.
func (m *MockEmailRepository) GetEmailByID(id int) (models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailByID", id)
	ret0, _ := ret[0].(models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmailByID indicates an expected call of GetEmailByID.
func (mr *MockEmailRepositoryMockRecorder) GetEmailByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailByID", reflect.TypeOf((*MockEmailRepository)(nil).GetEmailByID), id)
}

// GetSentEmails mocks base method.
func (m *MockEmailRepository) GetSentEmails(senderEmail string) ([]models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSentEmails", senderEmail)
	ret0, _ := ret[0].([]models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSentEmails indicates an expected call of GetSentEmails.
func (mr *MockEmailRepositoryMockRecorder) GetSentEmails(senderEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSentEmails", reflect.TypeOf((*MockEmailRepository)(nil).GetSentEmails), senderEmail)
}

// Inbox mocks base method.
func (m *MockEmailRepository) Inbox(id string) ([]models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inbox", id)
	ret0, _ := ret[0].([]models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inbox indicates an expected call of Inbox.
func (mr *MockEmailRepositoryMockRecorder) Inbox(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inbox", reflect.TypeOf((*MockEmailRepository)(nil).Inbox), id)
}

// SaveEmail mocks base method.
func (m *MockEmailRepository) SaveEmail(email models.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEmail", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEmail indicates an expected call of SaveEmail.
func (mr *MockEmailRepositoryMockRecorder) SaveEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEmail", reflect.TypeOf((*MockEmailRepository)(nil).SaveEmail), email)
}

// MockSMTPRepository is a mock of SMTPRepository interface.
type MockSMTPRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSMTPRepositoryMockRecorder
}

// MockSMTPRepositoryMockRecorder is the mock recorder for MockSMTPRepository.
type MockSMTPRepositoryMockRecorder struct {
	mock *MockSMTPRepository
}

// NewMockSMTPRepository creates a new mock instance.
func NewMockSMTPRepository(ctrl *gomock.Controller) *MockSMTPRepository {
	mock := &MockSMTPRepository{ctrl: ctrl}
	mock.recorder = &MockSMTPRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSMTPRepository) EXPECT() *MockSMTPRepositoryMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockSMTPRepository) SendEmail(from string, to []string, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", from, to, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockSMTPRepositoryMockRecorder) SendEmail(from, to, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockSMTPRepository)(nil).SendEmail), from, to, subject, body)
}

// MockPOP3Repository is a mock of POP3Repository interface.
type MockPOP3Repository struct {
	ctrl     *gomock.Controller
	recorder *MockPOP3RepositoryMockRecorder
}

// MockPOP3RepositoryMockRecorder is the mock recorder for MockPOP3Repository.
type MockPOP3RepositoryMockRecorder struct {
	mock *MockPOP3Repository
}

// NewMockPOP3Repository creates a new mock instance.
func NewMockPOP3Repository(ctrl *gomock.Controller) *MockPOP3Repository {
	mock := &MockPOP3Repository{ctrl: ctrl}
	mock.recorder = &MockPOP3RepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPOP3Repository) EXPECT() *MockPOP3RepositoryMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockPOP3Repository) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockPOP3RepositoryMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockPOP3Repository)(nil).Connect))
}

// FetchEmails mocks base method.
func (m *MockPOP3Repository) FetchEmails(arg0 models.EmailRepository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchEmails", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FetchEmails indicates an expected call of FetchEmails.
func (mr *MockPOP3RepositoryMockRecorder) FetchEmails(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchEmails", reflect.TypeOf((*MockPOP3Repository)(nil).FetchEmails), arg0)
}

// Quit mocks base method.
func (m *MockPOP3Repository) Quit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Quit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Quit indicates an expected call of Quit.
func (mr *MockPOP3RepositoryMockRecorder) Quit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Quit", reflect.TypeOf((*MockPOP3Repository)(nil).Quit))
}
