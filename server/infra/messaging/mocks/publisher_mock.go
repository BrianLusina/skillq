// Code generated by MockGen. DO NOT EDIT.
// Source: infra/messaging/publisher.go
//
// Generated by this command:
//
//	mockgen -source infra/messaging/publisher.go -destination infra/messaging/mocks/publisher_mock.go -package mockmessagepublisher
//

// Package mockmessagepublisher is a generated GoMock package.
package mockmessagepublisher

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPublisher is a mock of Publisher interface.
type MockPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockPublisherMockRecorder
}

// MockPublisherMockRecorder is the mock recorder for MockPublisher.
type MockPublisherMockRecorder struct {
	mock *MockPublisher
}

// NewMockPublisher creates a new mock instance.
func NewMockPublisher(ctrl *gomock.Controller) *MockPublisher {
	mock := &MockPublisher{ctrl: ctrl}
	mock.recorder = &MockPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublisher) EXPECT() *MockPublisherMockRecorder {
	return m.recorder
}

// CloseChan mocks base method.
func (m *MockPublisher) CloseChan() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseChan")
}

// CloseChan indicates an expected call of CloseChan.
func (mr *MockPublisherMockRecorder) CloseChan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseChan", reflect.TypeOf((*MockPublisher)(nil).CloseChan))
}

// Publish mocks base method.
func (m *MockPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, body, contentType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPublisherMockRecorder) Publish(ctx, body, contentType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPublisher)(nil).Publish), ctx, body, contentType)
}
