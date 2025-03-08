// Code generated by MockGen. DO NOT EDIT.
// Source: product.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockIProductRepo is a mock of IProductRepo interface.
type MockIProductRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIProductRepoMockRecorder
}

// MockIProductRepoMockRecorder is the mock recorder for MockIProductRepo.
type MockIProductRepoMockRecorder struct {
	mock *MockIProductRepo
}

// NewMockIProductRepo creates a new mock instance.
func NewMockIProductRepo(ctrl *gomock.Controller) *MockIProductRepo {
	mock := &MockIProductRepo{ctrl: ctrl}
	mock.recorder = &MockIProductRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProductRepo) EXPECT() *MockIProductRepoMockRecorder {
	return m.recorder
}

// GetAllProducts mocks base method.
func (m *MockIProductRepo) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProducts", ctx)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProducts indicates an expected call of GetAllProducts.
func (mr *MockIProductRepoMockRecorder) GetAllProducts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProducts", reflect.TypeOf((*MockIProductRepo)(nil).GetAllProducts), ctx)
}

// GetProductByID mocks base method.
func (m *MockIProductRepo) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByID", ctx, id)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID.
func (mr *MockIProductRepoMockRecorder) GetProductByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockIProductRepo)(nil).GetProductByID), ctx, id)
}

// GetProductCoverPath mocks base method.
func (m *MockIProductRepo) GetProductCoverPath(ctx context.Context, id int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductCoverPath", ctx, id)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductCoverPath indicates an expected call of GetProductCoverPath.
func (mr *MockIProductRepoMockRecorder) GetProductCoverPath(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductCoverPath", reflect.TypeOf((*MockIProductRepo)(nil).GetProductCoverPath), ctx, id)
}
