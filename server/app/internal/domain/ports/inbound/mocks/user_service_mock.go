// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/domain/ports/inbound/user_usecase.go

// Package mockusersvc is a generated GoMock package.
package mockusersvc

import (
	context "context"
	reflect "reflect"

	user "github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	inbound "github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	id "github.com/BrianLusina/skillq/server/domain/id"
	gomock "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// CreateEmailVerification mocks base method.
func (m *MockUserUseCase) CreateEmailVerification(arg0 context.Context, arg1 id.UUID) (user.UserVerification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmailVerification", arg0, arg1)
	ret0, _ := ret[0].(user.UserVerification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEmailVerification indicates an expected call of CreateEmailVerification.
func (mr *MockUserUseCaseMockRecorder) CreateEmailVerification(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmailVerification", reflect.TypeOf((*MockUserUseCase)(nil).CreateEmailVerification), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockUserUseCase) CreateUser(arg0 context.Context, arg1 inbound.UserRequest) (*inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserUseCaseMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserUseCase)(nil).CreateUser), arg0, arg1)
}

// GetUserByUUID mocks base method.
func (m *MockUserUseCase) GetUserByUUID(arg0 context.Context, arg1 id.UUID) (*inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUUID", arg0, arg1)
	ret0, _ := ret[0].(*inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUUID indicates an expected call of GetUserByUUID.
func (mr *MockUserUseCaseMockRecorder) GetUserByUUID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUUID", reflect.TypeOf((*MockUserUseCase)(nil).GetUserByUUID), arg0, arg1)
}

// UploadUserImage mocks base method.
func (m *MockUserUseCase) UploadUserImage(arg0 context.Context, arg1 id.UUID, arg2 inbound.UserImageRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadUserImage", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadUserImage indicates an expected call of UploadUserImage.
func (mr *MockUserUseCaseMockRecorder) UploadUserImage(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadUserImage", reflect.TypeOf((*MockUserUseCase)(nil).UploadUserImage), arg0, arg1, arg2)
}