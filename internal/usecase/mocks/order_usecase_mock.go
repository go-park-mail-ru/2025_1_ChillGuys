// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIOrderUsecase is a mock of IOrderUsecase interface.
type MockIOrderUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIOrderUsecaseMockRecorder
}

// MockIOrderUsecaseMockRecorder is the mock recorder for MockIOrderUsecase.
type MockIOrderUsecaseMockRecorder struct {
	mock *MockIOrderUsecase
}

// NewMockIOrderUsecase creates a new mock instance.
func NewMockIOrderUsecase(ctrl *gomock.Controller) *MockIOrderUsecase {
	mock := &MockIOrderUsecase{ctrl: ctrl}
	mock.recorder = &MockIOrderUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOrderUsecase) EXPECT() *MockIOrderUsecaseMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockIOrderUsecase) CreateOrder(arg0 context.Context, arg1 dto.CreateOrderDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockIOrderUsecaseMockRecorder) CreateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockIOrderUsecase)(nil).CreateOrder), arg0, arg1)
}

// GetOrdersPlaced mocks base method.
func (m *MockIOrderUsecase) GetOrdersPlaced(ctx context.Context) (*[]dto.OrderPreviewDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersPlaced", ctx)
	ret0, _ := ret[0].(*[]dto.OrderPreviewDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersPlaced indicates an expected call of GetOrdersPlaced.
func (mr *MockIOrderUsecaseMockRecorder) GetOrdersPlaced(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersPlaced", reflect.TypeOf((*MockIOrderUsecase)(nil).GetOrdersPlaced), ctx)
}

// GetUserOrders mocks base method.
func (m *MockIOrderUsecase) GetUserOrders(arg0 context.Context, arg1 uuid.UUID) (*[]dto.OrderPreviewDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", arg0, arg1)
	ret0, _ := ret[0].(*[]dto.OrderPreviewDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders.
func (mr *MockIOrderUsecaseMockRecorder) GetUserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockIOrderUsecase)(nil).GetUserOrders), arg0, arg1)
}

// UpdateStatus mocks base method.
func (m *MockIOrderUsecase) UpdateStatus(ctx context.Context, req dto.UpdateOrderStatusRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockIOrderUsecaseMockRecorder) UpdateStatus(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockIOrderUsecase)(nil).UpdateStatus), ctx, req)
}
