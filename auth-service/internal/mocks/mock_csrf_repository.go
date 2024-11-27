package mocks

import (
	"context"
	"mail/models"
	"reflect"

	"github.com/golang/mock/gomock"
)

type MockCsrfRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCsrfRepoMockRecorder
}

type MockCsrfRepoMockRecorder struct {
	mock *MockCsrfRepo
}

func NewMockCsrfRepo(ctrl *gomock.Controller) *MockCsrfRepo {
	mock := &MockCsrfRepo{ctrl: ctrl}
	mock.recorder = &MockCsrfRepoMockRecorder{mock}
	return mock
}

func (m *MockCsrfRepo) EXPECT() *MockCsrfRepoMockRecorder {
	return m.recorder
}

func (m *MockCsrfRepoMockRecorder) CreateCsrf(ctx, email interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "CreateCsrf", reflect.TypeOf((*MockCsrfRepo)(nil).CreateCsrf), ctx, email)
}

func (m *MockCsrfRepoMockRecorder) GetCsrf(ctx, token interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "GetCsrf", reflect.TypeOf((*MockCsrfRepo)(nil).GetCsrf), ctx, token)
}

func (m *MockCsrfRepoMockRecorder) DeleteCsrf(ctx, email interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "DeleteCsrf", reflect.TypeOf((*MockCsrfRepo)(nil).DeleteCsrf), ctx, email)
}

func (m *MockCsrfRepo) CreateCsrf(ctx context.Context, email string) (*models.Csrf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCsrf", ctx, email)
	ret0, _ := ret[0].(*models.Csrf)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockCsrfRepo) GetCsrf(ctx context.Context, token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCsrf", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockCsrfRepo) DeleteCsrf(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCsrf", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}
