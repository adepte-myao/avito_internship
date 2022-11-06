// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	dtos "github.com/adepte-myao/avito_internship/internal/dtos"
	models "github.com/adepte-myao/avito_internship/internal/models"
	gomock "github.com/golang/mock/gomock"
	decimal "github.com/shopspring/decimal"
)

// MockAccount is a mock of Account interface.
type MockAccount struct {
	ctrl     *gomock.Controller
	recorder *MockAccountMockRecorder
}

// MockAccountMockRecorder is the mock recorder for MockAccount.
type MockAccountMockRecorder struct {
	mock *MockAccount
}

// NewMockAccount creates a new mock instance.
func NewMockAccount(ctrl *gomock.Controller) *MockAccount {
	mock := &MockAccount{ctrl: ctrl}
	mock.recorder = &MockAccountMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccount) EXPECT() *MockAccountMockRecorder {
	return m.recorder
}

// Deposit mocks base method.
func (m *MockAccount) Deposit(accountId int32, value decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deposit", accountId, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deposit indicates an expected call of Deposit.
func (mr *MockAccountMockRecorder) Deposit(accountId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deposit", reflect.TypeOf((*MockAccount)(nil).Deposit), accountId, value)
}

// GetBalance mocks base method.
func (m *MockAccount) GetBalance(accountId int32) (models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", accountId)
	ret0, _ := ret[0].(models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockAccountMockRecorder) GetBalance(accountId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockAccount)(nil).GetBalance), accountId)
}

// Withdraw mocks base method.
func (m *MockAccount) Withdraw(accountId int32, value decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", accountId, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockAccountMockRecorder) Withdraw(accountId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockAccount)(nil).Withdraw), accountId, value)
}

// MockReservation is a mock of Reservation interface.
type MockReservation struct {
	ctrl     *gomock.Controller
	recorder *MockReservationMockRecorder
}

// MockReservationMockRecorder is the mock recorder for MockReservation.
type MockReservationMockRecorder struct {
	mock *MockReservation
}

// NewMockReservation creates a new mock instance.
func NewMockReservation(ctrl *gomock.Controller) *MockReservation {
	mock := &MockReservation{ctrl: ctrl}
	mock.recorder = &MockReservationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReservation) EXPECT() *MockReservationMockRecorder {
	return m.recorder
}

// AcceptReservation mocks base method.
func (m *MockReservation) AcceptReservation(resDto dtos.ReservationDto) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptReservation", resDto)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcceptReservation indicates an expected call of AcceptReservation.
func (mr *MockReservationMockRecorder) AcceptReservation(resDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptReservation", reflect.TypeOf((*MockReservation)(nil).AcceptReservation), resDto)
}

// CancelReservation mocks base method.
func (m *MockReservation) CancelReservation(resDto dtos.ReservationDto) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelReservation", resDto)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelReservation indicates an expected call of CancelReservation.
func (mr *MockReservationMockRecorder) CancelReservation(resDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelReservation", reflect.TypeOf((*MockReservation)(nil).CancelReservation), resDto)
}

// MakeReservation mocks base method.
func (m *MockReservation) MakeReservation(resDto dtos.ReservationDto) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeReservation", resDto)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeReservation indicates an expected call of MakeReservation.
func (mr *MockReservationMockRecorder) MakeReservation(resDto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeReservation", reflect.TypeOf((*MockReservation)(nil).MakeReservation), resDto)
}