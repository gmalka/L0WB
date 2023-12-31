// Code generated by MockGen. DO NOT EDIT.
// Source: store/cash/cash.go

package mock_orderservice

import (
	models "l0wb/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCasher is a mock of Casher interface.
type MockCasher struct {
	ctrl     *gomock.Controller
	recorder *MockCasherMockRecorder
}

// MockCasherMockRecorder is the mock recorder for MockCasher.
type MockCasherMockRecorder struct {
	mock *MockCasher
}

// NewMockCasher creates a new mock instance.
func NewMockCasher(ctrl *gomock.Controller) *MockCasher {
	mock := &MockCasher{ctrl: ctrl}
	mock.recorder = &MockCasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCasher) EXPECT() *MockCasherMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCasher) Add(arg0 models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockCasherMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCasher)(nil).Add), arg0)
}

// Get mocks base method.
func (m *MockCasher) Get(OrderUID string) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", OrderUID)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCasherMockRecorder) Get(OrderUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCasher)(nil).Get), OrderUID)
}
