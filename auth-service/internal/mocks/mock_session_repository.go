package mocks

import (
	"context"
	"mail/auth-service/internal/models"
	"reflect"

	"github.com/golang/mock/gomock"
)

type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
}

type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

func (m *MockSessionRepositoryMockRecorder) CreateSession(ctx, email interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "CreateSession", reflect.TypeOf((*MockSessionRepository)(nil).CreateSession), ctx, email)
}

func (m *MockSessionRepositoryMockRecorder) GetSession(ctx, id interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "GetSession", reflect.TypeOf((*MockSessionRepository)(nil).GetSession), ctx, id)
}

func (m *MockSessionRepositoryMockRecorder) DeleteSession(ctx, email interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "DeleteSession", reflect.TypeOf((*MockSessionRepository)(nil).DeleteSession), ctx, email)
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, email string) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, email)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockSessionRepository) GetSession(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockSessionRepository) DeleteSession(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}
