// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_cas is a generated GoMock package.
package mock_cas

import (
	cas "github.com/andrebq/cas"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockKV is a mock of KV interface
type MockKV struct {
	ctrl     *gomock.Controller
	recorder *MockKVMockRecorder
}

// MockKVMockRecorder is the mock recorder for MockKV
type MockKVMockRecorder struct {
	mock *MockKV
}

// NewMockKV creates a new mock instance
func NewMockKV(ctrl *gomock.Controller) *MockKV {
	mock := &MockKV{ctrl: ctrl}
	mock.recorder = &MockKVMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKV) EXPECT() *MockKVMockRecorder {
	return m.recorder
}

// Put mocks base method
func (m *MockKV) Put(k []byte, sz int) io.WriteCloser {
	ret := m.ctrl.Call(m, "Put", k, sz)
	ret0, _ := ret[0].(io.WriteCloser)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockKVMockRecorder) Put(k, sz interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockKV)(nil).Put), k, sz)
}

// Get mocks base method
func (m *MockKV) Get(k []byte) io.ReadCloser {
	ret := m.ctrl.Call(m, "Get", k)
	ret0, _ := ret[0].(io.ReadCloser)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockKVMockRecorder) Get(k interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKV)(nil).Get), k)
}

// MockOption is a mock of Option interface
type MockOption struct {
	ctrl     *gomock.Controller
	recorder *MockOptionMockRecorder
}

// MockOptionMockRecorder is the mock recorder for MockOption
type MockOptionMockRecorder struct {
	mock *MockOption
}

// NewMockOption creates a new mock instance
func NewMockOption(ctrl *gomock.Controller) *MockOption {
	mock := &MockOption{ctrl: ctrl}
	mock.recorder = &MockOptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOption) EXPECT() *MockOptionMockRecorder {
	return m.recorder
}

// apply mocks base method
func (m *MockOption) apply(arg0 *cas.Store) error {
	ret := m.ctrl.Call(m, "apply", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// apply indicates an expected call of apply
func (mr *MockOptionMockRecorder) apply(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "apply", reflect.TypeOf((*MockOption)(nil).apply), arg0)
}
