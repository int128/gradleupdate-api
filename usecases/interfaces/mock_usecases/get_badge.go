// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/int128/gradleupdate/usecases/interfaces (interfaces: GetBadge)

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/int128/gradleupdate/domain"
	interfaces "github.com/int128/gradleupdate/usecases/interfaces"
	reflect "reflect"
)

// MockGetBadge is a mock of GetBadge interface
type MockGetBadge struct {
	ctrl     *gomock.Controller
	recorder *MockGetBadgeMockRecorder
}

// MockGetBadgeMockRecorder is the mock recorder for MockGetBadge
type MockGetBadgeMockRecorder struct {
	mock *MockGetBadge
}

// NewMockGetBadge creates a new mock instance
func NewMockGetBadge(ctrl *gomock.Controller) *MockGetBadge {
	mock := &MockGetBadge{ctrl: ctrl}
	mock.recorder = &MockGetBadgeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGetBadge) EXPECT() *MockGetBadgeMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockGetBadge) Do(arg0 context.Context, arg1 domain.RepositoryID) (*interfaces.GetBadgeResponse, error) {
	ret := m.ctrl.Call(m, "Do", arg0, arg1)
	ret0, _ := ret[0].(*interfaces.GetBadgeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockGetBadgeMockRecorder) Do(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockGetBadge)(nil).Do), arg0, arg1)
}