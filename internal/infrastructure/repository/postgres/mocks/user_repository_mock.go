// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
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

// CheckExistence mocks base method.
func (m *MockIUserRepository) CheckExistence(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckExistence", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckExistence indicates an expected call of CheckExistence.
func (mr *MockIUserRepositoryMockRecorder) CheckExistence(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckExistence", reflect.TypeOf((*MockIUserRepository)(nil).CheckExistence), arg0, arg1)
}

// CheckVersion mocks base method.
func (m *MockIUserRepository) CheckVersion(arg0 context.Context, arg1 string, arg2 int) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckVersion", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckVersion indicates an expected call of CheckVersion.
func (mr *MockIUserRepositoryMockRecorder) CheckVersion(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckVersion", reflect.TypeOf((*MockIUserRepository)(nil).CheckVersion), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockIUserRepository) Create(arg0 context.Context, arg1 dto.UserDB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIUserRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIUserRepository)(nil).Create), arg0, arg1)
}

// GetByEmail mocks base method.
func (m *MockIUserRepository) GetByEmail(arg0 context.Context, arg1 string) (*dto.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", arg0, arg1)
	ret0, _ := ret[0].(*dto.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockIUserRepositoryMockRecorder) GetByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockIUserRepository)(nil).GetByEmail), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockIUserRepository) GetByID(arg0 context.Context, arg1 uuid.UUID) (*dto.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*dto.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIUserRepositoryMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIUserRepository)(nil).GetByID), arg0, arg1)
}

// GetCurrentVersion mocks base method.
func (m *MockIUserRepository) GetCurrentVersion(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentVersion", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentVersion indicates an expected call of GetCurrentVersion.
func (mr *MockIUserRepositoryMockRecorder) GetCurrentVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentVersion", reflect.TypeOf((*MockIUserRepository)(nil).GetCurrentVersion), arg0, arg1)
}

// IncrementVersion mocks base method.
func (m *MockIUserRepository) IncrementVersion(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementVersion", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementVersion indicates an expected call of IncrementVersion.
func (mr *MockIUserRepositoryMockRecorder) IncrementVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementVersion", reflect.TypeOf((*MockIUserRepository)(nil).IncrementVersion), arg0, arg1)
}

// UpdateImageURL mocks base method.
func (m *MockIUserRepository) UpdateImageURL(arg0 context.Context, arg1 uuid.UUID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateImageURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateImageURL indicates an expected call of UpdateImageURL.
func (mr *MockIUserRepositoryMockRecorder) UpdateImageURL(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateImageURL", reflect.TypeOf((*MockIUserRepository)(nil).UpdateImageURL), arg0, arg1, arg2)
}
