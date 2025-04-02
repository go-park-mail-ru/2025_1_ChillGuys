// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	jwt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockITokenator is a mock of ITokenator interface.
type MockITokenator struct {
	ctrl     *gomock.Controller
	recorder *MockITokenatorMockRecorder
}

// MockITokenatorMockRecorder is the mock recorder for MockITokenator.
type MockITokenatorMockRecorder struct {
	mock *MockITokenator
}

// NewMockITokenator creates a new mock instance.
func NewMockITokenator(ctrl *gomock.Controller) *MockITokenator {
	mock := &MockITokenator{ctrl: ctrl}
	mock.recorder = &MockITokenatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITokenator) EXPECT() *MockITokenatorMockRecorder {
	return m.recorder
}

// CreateJWT mocks base method.
func (m *MockITokenator) CreateJWT(userID string, version int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJWT", userID, version)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJWT indicates an expected call of CreateJWT.
func (mr *MockITokenatorMockRecorder) CreateJWT(userID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJWT", reflect.TypeOf((*MockITokenator)(nil).CreateJWT), userID, version)
}

// ParseJWT mocks base method.
func (m *MockITokenator) ParseJWT(tokenString string) (*jwt.JWTClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseJWT", tokenString)
	ret0, _ := ret[0].(*jwt.JWTClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseJWT indicates an expected call of ParseJWT.
func (mr *MockITokenatorMockRecorder) ParseJWT(tokenString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseJWT", reflect.TypeOf((*MockITokenator)(nil).ParseJWT), tokenString)
}

// MockIUserRepository is a mock of IUserRepository interface.
type MockIUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepositoryMockRecorder
}

// MockIUserRepositoryMockRecorder is the mock recorder for MockIUserRepository.
type MockIUserRepositoryMockRecorder struct {
	mock *MockIUserRepository
}

// NewMockIUserRepository creates a new mock instance.
func NewMockIUserRepository(ctrl *gomock.Controller) *MockIUserRepository {
	mock := &MockIUserRepository{ctrl: ctrl}
	mock.recorder = &MockIUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepository) EXPECT() *MockIUserRepositoryMockRecorder {
	return m.recorder
}

// CheckUserExists mocks base method.
func (m *MockIUserRepository) CheckUserExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserExists indicates an expected call of CheckUserExists.
func (mr *MockIUserRepositoryMockRecorder) CheckUserExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserExists", reflect.TypeOf((*MockIUserRepository)(nil).CheckUserExists), arg0, arg1)
}

// CheckUserVersion mocks base method.
func (m *MockIUserRepository) CheckUserVersion(arg0 context.Context, arg1 string, arg2 int) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserVersion", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckUserVersion indicates an expected call of CheckUserVersion.
func (mr *MockIUserRepositoryMockRecorder) CheckUserVersion(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserVersion", reflect.TypeOf((*MockIUserRepository)(nil).CheckUserVersion), arg0, arg1, arg2)
}

// CreateUser mocks base method.
func (m *MockIUserRepository) CreateUser(arg0 context.Context, arg1 models.UserDB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIUserRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIUserRepository)(nil).CreateUser), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockIUserRepository) GetUserByEmail(arg0 context.Context, arg1 string) (*models.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*models.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockIUserRepositoryMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockIUserRepository) GetUserByID(arg0 context.Context, arg1 uuid.UUID) (*models.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(*models.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockIUserRepositoryMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByID), arg0, arg1)
}

// GetUserCurrentVersion mocks base method.
func (m *MockIUserRepository) GetUserCurrentVersion(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCurrentVersion", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCurrentVersion indicates an expected call of GetUserCurrentVersion.
func (mr *MockIUserRepositoryMockRecorder) GetUserCurrentVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCurrentVersion", reflect.TypeOf((*MockIUserRepository)(nil).GetUserCurrentVersion), arg0, arg1)
}

// IncrementUserVersion mocks base method.
func (m *MockIUserRepository) IncrementUserVersion(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementUserVersion", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementUserVersion indicates an expected call of IncrementUserVersion.
func (mr *MockIUserRepositoryMockRecorder) IncrementUserVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementUserVersion", reflect.TypeOf((*MockIUserRepository)(nil).IncrementUserVersion), arg0, arg1)
}

// UpdateUserImageURL mocks base method.
func (m *MockIUserRepository) UpdateUserImageURL(arg0 context.Context, arg1 uuid.UUID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserImageURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserImageURL indicates an expected call of UpdateUserImageURL.
func (mr *MockIUserRepositoryMockRecorder) UpdateUserImageURL(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserImageURL", reflect.TypeOf((*MockIUserRepository)(nil).UpdateUserImageURL), arg0, arg1, arg2)
}
