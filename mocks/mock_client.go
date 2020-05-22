// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jaswanth05rongali/pub-sub/client (interfaces: Interface)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Init mocks base method.
func (m *MockInterface) Init() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init")
}

// Init indicates an expected call of Init.
func (mr *MockInterfaceMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockInterface)(nil).Init))
}

// RetrySendingMessage mocks base method.
func (m *MockInterface) RetrySendingMessage(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrySendingMessage", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// RetrySendingMessage indicates an expected call of RetrySendingMessage.
func (mr *MockInterfaceMockRecorder) RetrySendingMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrySendingMessage", reflect.TypeOf((*MockInterface)(nil).RetrySendingMessage), arg0)
}

// SaveToFile mocks base method.
func (m *MockInterface) SaveToFile(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveToFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveToFile indicates an expected call of SaveToFile.
func (mr *MockInterfaceMockRecorder) SaveToFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveToFile", reflect.TypeOf((*MockInterface)(nil).SaveToFile), arg0)
}

// SendMessage mocks base method.
func (m *MockInterface) SendMessage(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockInterfaceMockRecorder) SendMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockInterface)(nil).SendMessage), arg0)
}
