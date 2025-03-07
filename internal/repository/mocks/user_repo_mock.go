// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	jwt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

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

// CreateUser mocks base method.
func (m *MockIUserRepository) CreateUser(user models.UserRepo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIUserRepositoryMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIUserRepository)(nil).CreateUser), user)
}

// GetUserByEmail mocks base method.
func (m *MockIUserRepository) GetUserByEmail(email string) (*models.UserRepo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*models.UserRepo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockIUserRepositoryMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByEmail), email)
}

// GetUserByID mocks base method.
func (m *MockIUserRepository) GetUserByID(id uuid.UUID) (*models.UserRepo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(*models.UserRepo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockIUserRepositoryMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByID), id)
}

// IncrementUserVersion mocks base method.
func (m *MockIUserRepository) IncrementUserVersion(userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementUserVersion", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementUserVersion indicates an expected call of IncrementUserVersion.
func (mr *MockIUserRepositoryMockRecorder) IncrementUserVersion(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementUserVersion", reflect.TypeOf((*MockIUserRepository)(nil).IncrementUserVersion), userID)
}

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