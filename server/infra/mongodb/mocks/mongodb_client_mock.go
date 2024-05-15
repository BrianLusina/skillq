// Code generated by MockGen. DO NOT EDIT.
// Source: infra/mongodb/client.go
//
// Generated by this command:
//
//	mockgen -source infra/mongodb/client.go -destination infra/mongodb/mocks/mongodb_client_mock.go -package mockmongodb
//

// Package mockmongodb is a generated GoMock package.
package mockmongodb

import (
	context "context"
	reflect "reflect"

	mongodb "github.com/BrianLusina/skillq/server/infra/mongodb"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

// MockMongoDBClient is a mock of MongoDBClient interface.
type MockMongoDBClient[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockMongoDBClientMockRecorder[T]
}

// MockMongoDBClientMockRecorder is the mock recorder for MockMongoDBClient.
type MockMongoDBClientMockRecorder[T any] struct {
	mock *MockMongoDBClient[T]
}

// NewMockMongoDBClient creates a new mock instance.
func NewMockMongoDBClient[T any](ctrl *gomock.Controller) *MockMongoDBClient[T] {
	mock := &MockMongoDBClient[T]{ctrl: ctrl}
	mock.recorder = &MockMongoDBClientMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMongoDBClient[T]) EXPECT() *MockMongoDBClientMockRecorder[T] {
	return m.recorder
}

// BulkInsert mocks base method.
func (m *MockMongoDBClient[T]) BulkInsert(ctx context.Context, models []any) ([]primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkInsert", ctx, models)
	ret0, _ := ret[0].([]primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkInsert indicates an expected call of BulkInsert.
func (mr *MockMongoDBClientMockRecorder[T]) BulkInsert(ctx, models any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsert", reflect.TypeOf((*MockMongoDBClient[T])(nil).BulkInsert), ctx, models)
}

// Delete mocks base method.
func (m *MockMongoDBClient[T]) Delete(ctx context.Context, keyName, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, keyName, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockMongoDBClientMockRecorder[T]) Delete(ctx, keyName, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMongoDBClient[T])(nil).Delete), ctx, keyName, id)
}

// Disconnect mocks base method.
func (m *MockMongoDBClient[T]) Disconnect(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockMongoDBClientMockRecorder[T]) Disconnect(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockMongoDBClient[T])(nil).Disconnect), arg0)
}

// FindAll mocks base method.
func (m *MockMongoDBClient[T]) FindAll(ctx context.Context, filterOptions mongodb.FilterOptions) ([]T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, filterOptions)
	ret0, _ := ret[0].([]T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockMongoDBClientMockRecorder[T]) FindAll(ctx, filterOptions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockMongoDBClient[T])(nil).FindAll), ctx, filterOptions)
}

// FindById mocks base method.
func (m *MockMongoDBClient[T]) FindById(ctx context.Context, keyName, id string) (T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, keyName, id)
	ret0, _ := ret[0].(T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockMongoDBClientMockRecorder[T]) FindById(ctx, keyName, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockMongoDBClient[T])(nil).FindById), ctx, keyName, id)
}

// Insert mocks base method.
func (m *MockMongoDBClient[T]) Insert(ctx context.Context, model T) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, model)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockMongoDBClientMockRecorder[T]) Insert(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockMongoDBClient[T])(nil).Insert), ctx, model)
}
