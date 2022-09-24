// Code generated by MockGen. DO NOT EDIT.
// Source: backend.cs3219.comp.nus.edu.sg/auth (interfaces: TokenAuthenticator)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTokenAuthenticator is a mock of TokenAuthenticator interface.
type MockTokenAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockTokenAuthenticatorMockRecorder
}

// MockTokenAuthenticatorMockRecorder is the mock recorder for MockTokenAuthenticator.
type MockTokenAuthenticatorMockRecorder struct {
	mock *MockTokenAuthenticator
}

// NewMockTokenAuthenticator creates a new mock instance.
func NewMockTokenAuthenticator(ctrl *gomock.Controller) *MockTokenAuthenticator {
	mock := &MockTokenAuthenticator{ctrl: ctrl}
	mock.recorder = &MockTokenAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenAuthenticator) EXPECT() *MockTokenAuthenticatorMockRecorder {
	return m.recorder
}

// IsValidToken mocks base method.
func (m *MockTokenAuthenticator) IsValidToken(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidToken", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValidToken indicates an expected call of IsValidToken.
func (mr *MockTokenAuthenticatorMockRecorder) IsValidToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidToken", reflect.TypeOf((*MockTokenAuthenticator)(nil).IsValidToken), arg0)
}
