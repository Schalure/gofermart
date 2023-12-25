// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Schalure/gofermart/internal/gofermart (interfaces: Storager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	storage "github.com/Schalure/gofermart/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// AddNewOrder mocks base method.
func (m *MockStorager) AddNewOrder(arg0 context.Context, arg1 storage.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewOrder indicates an expected call of AddNewOrder.
func (mr *MockStoragerMockRecorder) AddNewOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewOrder", reflect.TypeOf((*MockStorager)(nil).AddNewOrder), arg0, arg1)
}

// AddNewUser mocks base method.
func (m *MockStorager) AddNewUser(arg0 context.Context, arg1 storage.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewUser indicates an expected call of AddNewUser.
func (mr *MockStoragerMockRecorder) AddNewUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewUser", reflect.TypeOf((*MockStorager)(nil).AddNewUser), arg0, arg1)
}

// DeleteOrder mocks base method.
func (m *MockStorager) DeleteOrder(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockStoragerMockRecorder) DeleteOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockStorager)(nil).DeleteOrder), arg0, arg1)
}

// GetOrderByNumber mocks base method.
func (m *MockStorager) GetOrderByNumber(arg0 context.Context, arg1 string) (storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByNumber", arg0, arg1)
	ret0, _ := ret[0].(storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByNumber indicates an expected call of GetOrderByNumber.
func (mr *MockStoragerMockRecorder) GetOrderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByNumber", reflect.TypeOf((*MockStorager)(nil).GetOrderByNumber), arg0, arg1)
}

// GetOrdersByLogin mocks base method.
func (m *MockStorager) GetOrdersByLogin(arg0 context.Context, arg1 string) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByLogin", arg0, arg1)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByLogin indicates an expected call of GetOrdersByLogin.
func (mr *MockStoragerMockRecorder) GetOrdersByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByLogin", reflect.TypeOf((*MockStorager)(nil).GetOrdersByLogin), arg0, arg1)
}

// GetOrdersToUpdateStatus mocks base method.
func (m *MockStorager) GetOrdersToUpdateStatus(arg0 context.Context) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersToUpdateStatus", arg0)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersToUpdateStatus indicates an expected call of GetOrdersToUpdateStatus.
func (mr *MockStoragerMockRecorder) GetOrdersToUpdateStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersToUpdateStatus", reflect.TypeOf((*MockStorager)(nil).GetOrdersToUpdateStatus), arg0)
}

// GetPointWithdraws mocks base method.
func (m *MockStorager) GetPointWithdraws(arg0 context.Context, arg1 string) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPointWithdraws", arg0, arg1)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPointWithdraws indicates an expected call of GetPointWithdraws.
func (mr *MockStoragerMockRecorder) GetPointWithdraws(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPointWithdraws", reflect.TypeOf((*MockStorager)(nil).GetPointWithdraws), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockStorager) GetUserByLogin(arg0 context.Context, arg1 string) (storage.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(storage.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStoragerMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorager)(nil).GetUserByLogin), arg0, arg1)
}

// UpdateOrder mocks base method.
func (m *MockStorager) UpdateOrder(arg0 context.Context, arg1, arg2 string, arg3 storage.OrderStatus, arg4 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStoragerMockRecorder) UpdateOrder(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStorager)(nil).UpdateOrder), arg0, arg1, arg2, arg3, arg4)
}

// WithdrawPointsForOrder mocks base method.
func (m *MockStorager) WithdrawPointsForOrder(arg0 context.Context, arg1 string, arg2 float64, arg3 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawPointsForOrder", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithdrawPointsForOrder indicates an expected call of WithdrawPointsForOrder.
func (mr *MockStoragerMockRecorder) WithdrawPointsForOrder(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawPointsForOrder", reflect.TypeOf((*MockStorager)(nil).WithdrawPointsForOrder), arg0, arg1, arg2, arg3)
}
