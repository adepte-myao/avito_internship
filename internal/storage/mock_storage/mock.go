// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	sql "database/sql"
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

// CreateAccount mocks base method.
func (m *MockAccount) CreateAccount(tx *sql.Tx, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", tx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAccountMockRecorder) CreateAccount(tx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccount)(nil).CreateAccount), tx, id)
}

// DecreaseBalance mocks base method.
func (m *MockAccount) DecreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecreaseBalance", tx, id, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecreaseBalance indicates an expected call of DecreaseBalance.
func (mr *MockAccountMockRecorder) DecreaseBalance(tx, id, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecreaseBalance", reflect.TypeOf((*MockAccount)(nil).DecreaseBalance), tx, id, value)
}

// GetAccount mocks base method.
func (m *MockAccount) GetAccount(tx *sql.Tx, id int32) (models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", tx, id)
	ret0, _ := ret[0].(models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountMockRecorder) GetAccount(tx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccount)(nil).GetAccount), tx, id)
}

// IncreaseBalance mocks base method.
func (m *MockAccount) IncreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseBalance", tx, id, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncreaseBalance indicates an expected call of IncreaseBalance.
func (mr *MockAccountMockRecorder) IncreaseBalance(tx, id, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseBalance", reflect.TypeOf((*MockAccount)(nil).IncreaseBalance), tx, id, value)
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

// CreateReservation mocks base method.
func (m *MockReservation) CreateReservation(tx *sql.Tx, reservation models.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReservation", tx, reservation)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReservation indicates an expected call of CreateReservation.
func (mr *MockReservationMockRecorder) CreateReservation(tx, reservation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReservation", reflect.TypeOf((*MockReservation)(nil).CreateReservation), tx, reservation)
}

// GetAccountantReport mocks base method.
func (m *MockReservation) GetAccountantReport(tx *sql.Tx) ([]models.AccountantReportElem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountantReport", tx)
	ret0, _ := ret[0].([]models.AccountantReportElem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountantReport indicates an expected call of GetAccountantReport.
func (mr *MockReservationMockRecorder) GetAccountantReport(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountantReport", reflect.TypeOf((*MockReservation)(nil).GetAccountantReport), tx)
}

// GetReservation mocks base method.
func (m *MockReservation) GetReservation(tx *sql.Tx, reservationDto dtos.ReservationDto, state models.ReserveState) (models.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReservation", tx, reservationDto, state)
	ret0, _ := ret[0].(models.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReservation indicates an expected call of GetReservation.
func (mr *MockReservationMockRecorder) GetReservation(tx, reservationDto, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReservation", reflect.TypeOf((*MockReservation)(nil).GetReservation), tx, reservationDto, state)
}

// MockTransfer is a mock of Transfer interface.
type MockTransfer struct {
	ctrl     *gomock.Controller
	recorder *MockTransferMockRecorder
}

// MockTransferMockRecorder is the mock recorder for MockTransfer.
type MockTransferMockRecorder struct {
	mock *MockTransfer
}

// NewMockTransfer creates a new mock instance.
func NewMockTransfer(ctrl *gomock.Controller) *MockTransfer {
	mock := &MockTransfer{ctrl: ctrl}
	mock.recorder = &MockTransferMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransfer) EXPECT() *MockTransferMockRecorder {
	return m.recorder
}

// GetAccountStatements mocks base method.
func (m *MockTransfer) GetAccountStatements(tx *sql.Tx, dto dtos.GetAccountStatementDto) ([]models.StatementElem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountStatements", tx, dto)
	ret0, _ := ret[0].([]models.StatementElem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountStatements indicates an expected call of GetAccountStatements.
func (mr *MockTransferMockRecorder) GetAccountStatements(tx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountStatements", reflect.TypeOf((*MockTransfer)(nil).GetAccountStatements), tx, dto)
}

// RecordExternalTransfer mocks base method.
func (m *MockTransfer) RecordExternalTransfer(tx *sql.Tx, accId int32, ttype models.TransferType, amount decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordExternalTransfer", tx, accId, ttype, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordExternalTransfer indicates an expected call of RecordExternalTransfer.
func (mr *MockTransferMockRecorder) RecordExternalTransfer(tx, accId, ttype, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordExternalTransfer", reflect.TypeOf((*MockTransfer)(nil).RecordExternalTransfer), tx, accId, ttype, amount)
}

// RecordInternalTransfer mocks base method.
func (m *MockTransfer) RecordInternalTransfer(tx *sql.Tx, senderId, recId int32, amount decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordInternalTransfer", tx, senderId, recId, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordInternalTransfer indicates an expected call of RecordInternalTransfer.
func (mr *MockTransferMockRecorder) RecordInternalTransfer(tx, senderId, recId, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordInternalTransfer", reflect.TypeOf((*MockTransfer)(nil).RecordInternalTransfer), tx, senderId, recId, amount)
}

// MockSQLTransactionHelper is a mock of SQLTransactionHelper interface.
type MockSQLTransactionHelper struct {
	ctrl     *gomock.Controller
	recorder *MockSQLTransactionHelperMockRecorder
}

// MockSQLTransactionHelperMockRecorder is the mock recorder for MockSQLTransactionHelper.
type MockSQLTransactionHelperMockRecorder struct {
	mock *MockSQLTransactionHelper
}

// NewMockSQLTransactionHelper creates a new mock instance.
func NewMockSQLTransactionHelper(ctrl *gomock.Controller) *MockSQLTransactionHelper {
	mock := &MockSQLTransactionHelper{ctrl: ctrl}
	mock.recorder = &MockSQLTransactionHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSQLTransactionHelper) EXPECT() *MockSQLTransactionHelperMockRecorder {
	return m.recorder
}

// BeginTransaction mocks base method.
func (m *MockSQLTransactionHelper) BeginTransaction() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTransaction")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTransaction indicates an expected call of BeginTransaction.
func (mr *MockSQLTransactionHelperMockRecorder) BeginTransaction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTransaction", reflect.TypeOf((*MockSQLTransactionHelper)(nil).BeginTransaction))
}

// CommitTransaction mocks base method.
func (m *MockSQLTransactionHelper) CommitTransaction(tx *sql.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CommitTransaction", tx)
}

// CommitTransaction indicates an expected call of CommitTransaction.
func (mr *MockSQLTransactionHelperMockRecorder) CommitTransaction(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitTransaction", reflect.TypeOf((*MockSQLTransactionHelper)(nil).CommitTransaction), tx)
}

// RollbackTransaction mocks base method.
func (m *MockSQLTransactionHelper) RollbackTransaction(tx *sql.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RollbackTransaction", tx)
}

// RollbackTransaction indicates an expected call of RollbackTransaction.
func (mr *MockSQLTransactionHelperMockRecorder) RollbackTransaction(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTransaction", reflect.TypeOf((*MockSQLTransactionHelper)(nil).RollbackTransaction), tx)
}
