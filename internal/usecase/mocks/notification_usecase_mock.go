// Code generated by MockGen. DO NOT EDIT.
// Source: notification.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockINotificationUsecase is a mock of INotificationUsecase interface.
type MockINotificationUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockINotificationUsecaseMockRecorder
}

// MockINotificationUsecaseMockRecorder is the mock recorder for MockINotificationUsecase.
type MockINotificationUsecaseMockRecorder struct {
	mock *MockINotificationUsecase
}

// NewMockINotificationUsecase creates a new mock instance.
func NewMockINotificationUsecase(ctrl *gomock.Controller) *MockINotificationUsecase {
	mock := &MockINotificationUsecase{ctrl: ctrl}
	mock.recorder = &MockINotificationUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockINotificationUsecase) EXPECT() *MockINotificationUsecaseMockRecorder {
	return m.recorder
}

// GetAllByUser mocks base method.
func (m *MockINotificationUsecase) GetAllByUser(ctx context.Context, offset int) (dto.NotificationsListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUser", ctx, offset)
	ret0, _ := ret[0].(dto.NotificationsListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUser indicates an expected call of GetAllByUser.
func (mr *MockINotificationUsecaseMockRecorder) GetAllByUser(ctx, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUser", reflect.TypeOf((*MockINotificationUsecase)(nil).GetAllByUser), ctx, offset)
}

// GetUnreadCount mocks base method.
func (m *MockINotificationUsecase) GetUnreadCount(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadCount", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnreadCount indicates an expected call of GetUnreadCount.
func (mr *MockINotificationUsecaseMockRecorder) GetUnreadCount(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadCount", reflect.TypeOf((*MockINotificationUsecase)(nil).GetUnreadCount), ctx)
}

// MarkAsRead mocks base method.
func (m *MockINotificationUsecase) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsRead", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAsRead indicates an expected call of MarkAsRead.
func (mr *MockINotificationUsecaseMockRecorder) MarkAsRead(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsRead", reflect.TypeOf((*MockINotificationUsecase)(nil).MarkAsRead), ctx, id)
}
