// Code generated by MockGen. DO NOT EDIT.
// Source: search.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	dto "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gomock "github.com/golang/mock/gomock"
	null "github.com/guregu/null"
)

// MockISearchUsecase is a mock of ISearchUsecase interface.
type MockISearchUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockISearchUsecaseMockRecorder
}

// MockISearchUsecaseMockRecorder is the mock recorder for MockISearchUsecase.
type MockISearchUsecaseMockRecorder struct {
	mock *MockISearchUsecase
}

// NewMockISearchUsecase creates a new mock instance.
func NewMockISearchUsecase(ctrl *gomock.Controller) *MockISearchUsecase {
	mock := &MockISearchUsecase{ctrl: ctrl}
	mock.recorder = &MockISearchUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISearchUsecase) EXPECT() *MockISearchUsecaseMockRecorder {
	return m.recorder
}

// SearchCategoryByName mocks base method.
func (m *MockISearchUsecase) SearchCategoryByName(arg0 context.Context, arg1 dto.CategoryNameResponse) ([]*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCategoryByName", arg0, arg1)
	ret0, _ := ret[0].([]*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchCategoryByName indicates an expected call of SearchCategoryByName.
func (mr *MockISearchUsecaseMockRecorder) SearchCategoryByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCategoryByName", reflect.TypeOf((*MockISearchUsecase)(nil).SearchCategoryByName), arg0, arg1)
}

// SearchProductsByNameWithFilterAndSort mocks base method.
func (m *MockISearchUsecase) SearchProductsByNameWithFilterAndSort(ctx context.Context, categoryID null.String, subString string, offset int, minPrice, maxPrice float64, minRating float32, sortOption models.SortOption) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchProductsByNameWithFilterAndSort", ctx, categoryID, subString, offset, minPrice, maxPrice, minRating, sortOption)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchProductsByNameWithFilterAndSort indicates an expected call of SearchProductsByNameWithFilterAndSort.
func (mr *MockISearchUsecaseMockRecorder) SearchProductsByNameWithFilterAndSort(ctx, categoryID, subString, offset, minPrice, maxPrice, minRating, sortOption interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchProductsByNameWithFilterAndSort", reflect.TypeOf((*MockISearchUsecase)(nil).SearchProductsByNameWithFilterAndSort), ctx, categoryID, subString, offset, minPrice, maxPrice, minRating, sortOption)
}
