// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	minio "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	dto "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gomock "github.com/golang/mock/gomock"
)

// MockIUserUsecase is a mock of IUserUsecase interface.
type MockIUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIUserUsecaseMockRecorder
}

// MockIUserUsecaseMockRecorder is the mock recorder for MockIUserUsecase.
type MockIUserUsecaseMockRecorder struct {
	mock *MockIUserUsecase
}

// NewMockIUserUsecase creates a new mock instance.
func NewMockIUserUsecase(ctrl *gomock.Controller) *MockIUserUsecase {
	mock := &MockIUserUsecase{ctrl: ctrl}
	mock.recorder = &MockIUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserUsecase) EXPECT() *MockIUserUsecaseMockRecorder {
	return m.recorder
}

// GetMe mocks base method.
func (m *MockIUserUsecase) GetMe(arg0 context.Context) (*dto.UserDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMe", arg0)
	ret0, _ := ret[0].(*dto.UserDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockIUserUsecaseMockRecorder) GetMe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockIUserUsecase)(nil).GetMe), arg0)
}

// UpdateUserEmail mocks base method.
func (m *MockIUserUsecase) UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserEmail", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserEmail indicates an expected call of UpdateUserEmail.
func (mr *MockIUserUsecaseMockRecorder) UpdateUserEmail(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserEmail", reflect.TypeOf((*MockIUserUsecase)(nil).UpdateUserEmail), ctx, user)
}

// UpdateUserPassword mocks base method.
func (m *MockIUserUsecase) UpdateUserPassword(arg0 context.Context, arg1 dto.UpdateUserPasswordDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockIUserUsecaseMockRecorder) UpdateUserPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockIUserUsecase)(nil).UpdateUserPassword), arg0, arg1)
}

// UpdateUserProfile mocks base method.
func (m *MockIUserUsecase) UpdateUserProfile(arg0 context.Context, arg1 dto.UpdateUserProfileRequestDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockIUserUsecaseMockRecorder) UpdateUserProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockIUserUsecase)(nil).UpdateUserProfile), arg0, arg1)
}

// UploadAvatar mocks base method.
func (m *MockIUserUsecase) UploadAvatar(arg0 context.Context, arg1 minio.FileData) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadAvatar", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadAvatar indicates an expected call of UploadAvatar.
func (mr *MockIUserUsecaseMockRecorder) UploadAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadAvatar", reflect.TypeOf((*MockIUserUsecase)(nil).UploadAvatar), arg0, arg1)
}
