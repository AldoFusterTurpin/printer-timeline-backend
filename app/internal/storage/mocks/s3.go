// Code generated by MockGen. DO NOT EDIT.
// Source: ../storage/s3.go

// Package mocks is a generated GoMock package.
package mocks

import (
	s3 "github.com/aws/aws-sdk-go/service/s3"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockS3Fetcher is a mock of S3Fetcher interface
type MockS3Fetcher struct {
	ctrl     *gomock.Controller
	recorder *MockS3FetcherMockRecorder
}

// MockS3FetcherMockRecorder is the mock recorder for MockS3Fetcher
type MockS3FetcherMockRecorder struct {
	mock *MockS3Fetcher
}

// NewMockS3Fetcher creates a new mock instance
func NewMockS3Fetcher(ctrl *gomock.Controller) *MockS3Fetcher {
	mock := &MockS3Fetcher{ctrl: ctrl}
	mock.recorder = &MockS3FetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockS3Fetcher) EXPECT() *MockS3FetcherMockRecorder {
	return m.recorder
}

// GetObject mocks base method
func (m *MockS3Fetcher) GetObject(arg0 *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObject", arg0)
	ret0, _ := ret[0].(*s3.GetObjectOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObject indicates an expected call of GetObject
func (mr *MockS3FetcherMockRecorder) GetObject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObject", reflect.TypeOf((*MockS3Fetcher)(nil).GetObject), arg0)
}
