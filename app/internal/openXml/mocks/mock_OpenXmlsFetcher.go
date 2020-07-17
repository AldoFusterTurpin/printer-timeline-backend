// Code generated by MockGen. DO NOT EDIT.
// Source: bitbucket.org/aldoft/printer-timeline-backend/openXml (interfaces: OpenXmlsFetcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	cloudwatchlogs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockOpenXmlsFetcher is a mock of OpenXmlsFetcher interface
type MockOpenXmlsFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockOpenXmlsFetcherMockRecorder
}

// MockOpenXmlsFetcherMockRecorder is the mock recorder for MockOpenXmlsFetcher
type MockOpenXmlsFetcherMockRecorder struct {
	mock *MockOpenXmlsFetcher
}

// NewMockOpenXmlsFetcher creates a new mock instance
func NewMockOpenXmlsFetcher(ctrl *gomock.Controller) *MockOpenXmlsFetcher {
	mock := &MockOpenXmlsFetcher{ctrl: ctrl}
	mock.recorder = &MockOpenXmlsFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOpenXmlsFetcher) EXPECT() *MockOpenXmlsFetcherMockRecorder {
	return m.recorder
}

// GetUploadedOpenXmls mocks base method
func (m *MockOpenXmlsFetcher) GetUploadedOpenXmls(arg0 map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUploadedOpenXmls", arg0)
	ret0, _ := ret[0].(*cloudwatchlogs.GetQueryResultsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUploadedOpenXmls indicates an expected call of GetUploadedOpenXmls
func (mr *MockOpenXmlsFetcherMockRecorder) GetUploadedOpenXmls(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUploadedOpenXmls", reflect.TypeOf((*MockOpenXmlsFetcher)(nil).GetUploadedOpenXmls), arg0)
}