// Code generated by MockGen. DO NOT EDIT.
// Source: internal/jwt/jwt_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	jwt "github.com/MaxFando/lms/user-service/internal/jwt"
	model "github.com/MaxFando/lms/user-service/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockJWTService is a mock of JWTService interface.
type MockJWTService struct {
	ctrl     *gomock.Controller
	recorder *MockJWTServiceMockRecorder
}

// MockJWTServiceMockRecorder is the mock recorder for MockJWTService.
type MockJWTServiceMockRecorder struct {
	mock *MockJWTService
}

// NewMockJWTService creates a new mock instance.
func NewMockJWTService(ctrl *gomock.Controller) *MockJWTService {
	mock := &MockJWTService{ctrl: ctrl}
	mock.recorder = &MockJWTServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTService) EXPECT() *MockJWTServiceMockRecorder {
	return m.recorder
}

// GenerateTokens mocks base method.
func (m *MockJWTService) GenerateTokens(user *model.User) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokens", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateTokens indicates an expected call of GenerateTokens.
func (mr *MockJWTServiceMockRecorder) GenerateTokens(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokens", reflect.TypeOf((*MockJWTService)(nil).GenerateTokens), user)
}

// ParseToken mocks base method.
func (m *MockJWTService) ParseToken(token string) (*jwt.UserClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(*jwt.UserClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockJWTServiceMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockJWTService)(nil).ParseToken), token)
}
