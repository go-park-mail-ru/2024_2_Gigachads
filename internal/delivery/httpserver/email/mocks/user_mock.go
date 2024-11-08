// Code generated by MockGen. DO NOT EDIT.
// Source: internal/models/user_model.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "mail/internal/models"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// ChangeAvatar mocks base method.
func (m *MockUserUseCase) ChangeAvatar(file multipart.File, header multipart.FileHeader, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAvatar", file, header, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAvatar indicates an expected call of ChangeAvatar.
func (mr *MockUserUseCaseMockRecorder) ChangeAvatar(file, header, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAvatar", reflect.TypeOf((*MockUserUseCase)(nil).ChangeAvatar), file, header, email)
}

// ChangeName mocks base method.
func (m *MockUserUseCase) ChangeName(email, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeName", email, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeName indicates an expected call of ChangeName.
func (mr *MockUserUseCaseMockRecorder) ChangeName(email, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeName", reflect.TypeOf((*MockUserUseCase)(nil).ChangeName), email, name)
}

// ChangePassword mocks base method.
func (m *MockUserUseCase) ChangePassword(email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockUserUseCaseMockRecorder) ChangePassword(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockUserUseCase)(nil).ChangePassword), email, password)
}

// CheckAuth mocks base method.
func (m *MockUserUseCase) CheckAuth(ctx context.Context, sessionID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", ctx, sessionID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockUserUseCaseMockRecorder) CheckAuth(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockUserUseCase)(nil).CheckAuth), ctx, sessionID)
}

// CheckCsrf mocks base method.
func (m *MockUserUseCase) CheckCsrf(ctx context.Context, sessionID, scrf string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckCsrf", ctx, sessionID, scrf)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckCsrf indicates an expected call of CheckCsrf.
func (mr *MockUserUseCaseMockRecorder) CheckCsrf(ctx, sessionID, scrf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCsrf", reflect.TypeOf((*MockUserUseCase)(nil).CheckCsrf), ctx, sessionID, scrf)
}

// GetAvatar mocks base method.
func (m *MockUserUseCase) GetAvatar(email string) ([]byte, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", email)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockUserUseCaseMockRecorder) GetAvatar(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockUserUseCase)(nil).GetAvatar), email)
}

// Login mocks base method.
func (m *MockUserUseCase) Login(ctx context.Context, user *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*models.Session)
	ret2, _ := ret[2].(*models.Csrf)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Login indicates an expected call of Login.
func (mr *MockUserUseCaseMockRecorder) Login(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUseCase)(nil).Login), ctx, user)
}

// Logout mocks base method.
func (m *MockUserUseCase) Logout(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockUserUseCaseMockRecorder) Logout(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockUserUseCase)(nil).Logout), ctx, id)
}

// Signup mocks base method.
func (m *MockUserUseCase) Signup(ctx context.Context, user *models.User) (*models.User, *models.Session, *models.Csrf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*models.Session)
	ret2, _ := ret[2].(*models.Csrf)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Signup indicates an expected call of Signup.
func (mr *MockUserUseCaseMockRecorder) Signup(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockUserUseCase)(nil).Signup), ctx, user)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockUserRepository) CheckUser(login *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", login)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockUserRepositoryMockRecorder) CheckUser(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockUserRepository)(nil).CheckUser), login)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(signup *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", signup)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(signup interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), signup)
}

// GetUserByEmail mocks base method.
func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepositoryMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserByEmail), email)
}

// IsExist mocks base method.
func (m *MockUserRepository) IsExist(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExist", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExist indicates an expected call of IsExist.
func (mr *MockUserRepositoryMockRecorder) IsExist(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExist", reflect.TypeOf((*MockUserRepository)(nil).IsExist), email)
}

// UpdateInfo mocks base method.
func (m *MockUserRepository) UpdateInfo(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInfo", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInfo indicates an expected call of UpdateInfo.
func (mr *MockUserRepositoryMockRecorder) UpdateInfo(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInfo", reflect.TypeOf((*MockUserRepository)(nil).UpdateInfo), user)
}
