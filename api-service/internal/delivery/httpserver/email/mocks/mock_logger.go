// Code generated by MockGen. DO NOT EDIT.
// Source: api-service/pkg/logger/logger.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogable is a mock of Logable interface.
type MockLogable struct {
	ctrl     *gomock.Controller
	recorder *MockLogableMockRecorder
}

// MockLogableMockRecorder is the mock recorder for MockLogable.
type MockLogableMockRecorder struct {
	mock *MockLogable
}

// NewMockLogable creates a new mock instance.
func NewMockLogable(ctrl *gomock.Controller) *MockLogable {
	mock := &MockLogable{ctrl: ctrl}
	mock.recorder = &MockLogableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogable) EXPECT() *MockLogableMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogable) Debug(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLogableMockRecorder) Debug(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogable)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockLogable) Error(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLogableMockRecorder) Error(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogable)(nil).Error), varargs...)
}

// Info mocks base method.
func (m *MockLogable) Info(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLogableMockRecorder) Info(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogable)(nil).Info), varargs...)
}

// Warn mocks base method.
func (m *MockLogable) Warn(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockLogableMockRecorder) Warn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogable)(nil).Warn), varargs...)
}
