// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/domain/ports/inbound/user_service.go
//
// Generated by this command:
//
//	mockgen -source app/internal/domain/ports/inbound/user_service.go -destination app/internal/domain/ports/inbound/mocks/user_service_mock.go -package mockusersvc
//

// Package mockusersvc is a generated GoMock package.
package mockusersvc

import (
	context "context"
	reflect "reflect"

	inbound "github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	common "github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	id "github.com/BrianLusina/skillq/server/domain/id"
	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(arg0 context.Context, arg1 inbound.UserRequest) (*inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockUserService) DeleteUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserServiceMockRecorder) DeleteUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserService)(nil).DeleteUser), arg0, arg1)
}

// GetAllUsers mocks base method.
func (m *MockUserService) GetAllUsers(arg0 context.Context, arg1 common.RequestParams) ([]inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", arg0, arg1)
	ret0, _ := ret[0].([]inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockUserServiceMockRecorder) GetAllUsers(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUserService)(nil).GetAllUsers), arg0, arg1)
}

// GetAllUsersBySkill mocks base method.
func (m *MockUserService) GetAllUsersBySkill(arg0 context.Context, arg1 string, arg2 common.RequestParams) ([]inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsersBySkill", arg0, arg1, arg2)
	ret0, _ := ret[0].([]inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsersBySkill indicates an expected call of GetAllUsersBySkill.
func (mr *MockUserServiceMockRecorder) GetAllUsersBySkill(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsersBySkill", reflect.TypeOf((*MockUserService)(nil).GetAllUsersBySkill), arg0, arg1, arg2)
}

// GetUserByUUID mocks base method.
func (m *MockUserService) GetUserByUUID(arg0 context.Context, arg1 string) (*inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUUID", arg0, arg1)
	ret0, _ := ret[0].(*inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUUID indicates an expected call of GetUserByUUID.
func (mr *MockUserServiceMockRecorder) GetUserByUUID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUUID", reflect.TypeOf((*MockUserService)(nil).GetUserByUUID), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockUserService) UpdateUser(ctx context.Context, userID string, request inbound.UserRequest) (*inbound.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, userID, request)
	ret0, _ := ret[0].(*inbound.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserServiceMockRecorder) UpdateUser(ctx, userID, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserService)(nil).UpdateUser), ctx, userID, request)
}

// UploadUserImage mocks base method.
func (m *MockUserService) UploadUserImage(arg0 context.Context, arg1 id.UUID, arg2 inbound.UserImageRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadUserImage", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadUserImage indicates an expected call of UploadUserImage.
func (mr *MockUserServiceMockRecorder) UploadUserImage(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadUserImage", reflect.TypeOf((*MockUserService)(nil).UploadUserImage), arg0, arg1, arg2)
}
