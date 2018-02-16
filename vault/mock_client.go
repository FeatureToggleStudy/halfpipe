// Code generated by MockGen. DO NOT EDIT.
// Source: vault/client.go

// Package mock_vault is a generated GoMock package.
package vault

import (
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Exists mocks base method
func (m *MockClient) Exists(team, pipeline, mapKey, keyName string) (bool, error) {
	ret := m.ctrl.Call(m, "Exists", team, pipeline, mapKey, keyName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists
func (mr *MockClientMockRecorder) Exists(team, pipeline, mapKey, keyName interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockClient)(nil).Exists), team, pipeline, mapKey, keyName)
}

// VaultPrefix mocks base method
func (m *MockClient) VaultPrefix() string {
	ret := m.ctrl.Call(m, "VaultPrefix")
	ret0, _ := ret[0].(string)
	return ret0
}

// VaultPrefix indicates an expected call of VaultPrefix
func (mr *MockClientMockRecorder) VaultPrefix() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VaultPrefix", reflect.TypeOf((*MockClient)(nil).VaultPrefix))
}