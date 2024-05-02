// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/domain/ports/outbound/repositories/user_repo_port.go
//
// Generated by this command:
//
//	mockgen -source app/internal/domain/ports/outbound/repositories/user_repo_port.go -destination app/internal/domain/ports/outbound/repositories/mocks/user_repo_port_mock.go -package mockuserrepo
//

// Package mockuserrepo is a generated GoMock package.
package mockuserrepo

import (
	context "context"
	reflect "reflect"

	user "github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	id "github.com/BrianLusina/skillq/server/domain/id"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepoPort is a mock of UserRepoPort interface.
type MockUserRepoPort struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoPortMockRecorder
}

// MockUserRepoPortMockRecorder is the mock recorder for MockUserRepoPort.
type MockUserRepoPortMockRecorder struct {
	mock *MockUserRepoPort
}

// NewMockUserRepoPort creates a new mock instance.
func NewMockUserRepoPort(ctrl *gomock.Controller) *MockUserRepoPort {
	mock := &MockUserRepoPort{ctrl: ctrl}
	mock.recorder = &MockUserRepoPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepoPort) EXPECT() *MockUserRepoPortMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepoPort) CreateUser(arg0 context.Context, arg1 user.User) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoPortMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepoPort)(nil).CreateUser), arg0, arg1)
}

// GetAllUsers mocks base method.
func (m *MockUserRepoPort) GetAllUsers(arg0 context.Context) ([]user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", arg0)
	ret0, _ := ret[0].([]user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockUserRepoPortMockRecorder) GetAllUsers(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUserRepoPort)(nil).GetAllUsers), arg0)
}

// GetUserByUUID mocks base method.
func (m *MockUserRepoPort) GetUserByUUID(arg0 context.Context, arg1 id.UUID) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUUID", arg0, arg1)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUUID indicates an expected call of GetUserByUUID.
func (mr *MockUserRepoPortMockRecorder) GetUserByUUID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUUID", reflect.TypeOf((*MockUserRepoPort)(nil).GetUserByUUID), arg0, arg1)
}
