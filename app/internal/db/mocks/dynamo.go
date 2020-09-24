// Code generated by MockGen. DO NOT EDIT.
// Source: db/dynamo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	db "bitbucket.org/aldoft/printer-timeline-backend/app/internal/db"
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPrinterSubscriptionFetcher is a mock of PrinterSubscriptionFetcher interface
type MockPrinterSubscriptionFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockPrinterSubscriptionFetcherMockRecorder
}

// MockPrinterSubscriptionFetcherMockRecorder is the mock recorder for MockPrinterSubscriptionFetcher
type MockPrinterSubscriptionFetcherMockRecorder struct {
	mock *MockPrinterSubscriptionFetcher
}

// NewMockPrinterSubscriptionFetcher creates a new mock instance
func NewMockPrinterSubscriptionFetcher(ctrl *gomock.Controller) *MockPrinterSubscriptionFetcher {
	mock := &MockPrinterSubscriptionFetcher{ctrl: ctrl}
	mock.recorder = &MockPrinterSubscriptionFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPrinterSubscriptionFetcher) EXPECT() *MockPrinterSubscriptionFetcherMockRecorder {
	return m.recorder
}

// GetPrinterSubscriptions mocks base method
func (m *MockPrinterSubscriptionFetcher) GetPrinterSubscriptions(ctx context.Context, printerId string) ([]*db.CCPrinterSubscriptionModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrinterSubscriptions", ctx, printerId)
	ret0, _ := ret[0].([]*db.CCPrinterSubscriptionModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrinterSubscriptions indicates an expected call of GetPrinterSubscriptions
func (mr *MockPrinterSubscriptionFetcherMockRecorder) GetPrinterSubscriptions(ctx, printerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrinterSubscriptions", reflect.TypeOf((*MockPrinterSubscriptionFetcher)(nil).GetPrinterSubscriptions), ctx, printerId)
}