// Code generated by MockGen. DO NOT EDIT.
// Source: transport/resttransport/rest.go

// Package mock_resttransport is a generated GoMock package.
package mock_resttransport

import (
	models "l0wb/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOrderer is a mock of Orderer interface.
type MockOrderer struct {
	ctrl     *gomock.Controller
	recorder *MockOrdererMockRecorder
}

// MockOrdererMockRecorder is the mock recorder for MockOrderer.
type MockOrdererMockRecorder struct {
	mock *MockOrderer
}

// NewMockOrderer creates a new mock instance.
func NewMockOrderer(ctrl *gomock.Controller) *MockOrderer {
	mock := &MockOrderer{ctrl: ctrl}
	mock.recorder = &MockOrdererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderer) EXPECT() *MockOrdererMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockOrderer) Add(arg0 models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockOrdererMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockOrderer)(nil).Add), arg0)
}

// Get mocks base method.
func (m *MockOrderer) Get(OrderUID string) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", OrderUID)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockOrdererMockRecorder) Get(OrderUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockOrderer)(nil).Get), OrderUID)
}
