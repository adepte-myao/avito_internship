package services_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/adepte-myao/avito_internship/internal/storage/mock_storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestReservationer_MakeReservation(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal)
	type reservationRepoBehavior func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                    string
		inputDto                dtos.ReservationDto
		accountRepoBehavior     accountRepoBehavior
		reservationRepoBehavior reservationRepoBehavior
		txHelperBehavior        txHelperBehavior
		expectedError           error
	}{
		{
			name:     "Success",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromInt(100)}, nil,
				)
				accountRepo.EXPECT().DecreaseBalance(tx, accountId, value).Return(nil)
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().CreateReservation(tx, gomock.AssignableToTypeOf(models.Reservation{})).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
		{
			name:     "Fail account does not exist",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{}, errors.New("not nil"),
				)
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("account with ID 1 does not exist"),
		},
		{
			name:     "Fail not enough money",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromInt(99)}, nil,
				)
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("not enough money: account: id: 1, balance: 99; required: 100"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, testCase.inputDto.AccountId, testCase.inputDto.TotalCost)

			resRepo := mock_storage.NewMockReservation(ctrl)
			testCase.reservationRepoBehavior(resRepo, tx)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			reservationer := services.Reservationer{
				Account:     accRepo,
				Reservation: resRepo,
				TxHelper:    txHelper,
			}

			err := reservationer.MakeReservation(testCase.inputDto)
			assert.EqualValues(t, testCase.expectedError, err)
		})
	}
}

func TestReservationer_AcceptReservation(t *testing.T) {
	type reservationRepoBehavior func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                    string
		inputDto                dtos.ReservationDto
		reservationRepoBehavior reservationRepoBehavior
		txHelperBehavior        txHelperBehavior
		expectedError           error
	}{
		{
			name:     "Success",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Reserved).
					Return(models.Reservation{}, nil)
				reservationRepo.EXPECT().CreateReservation(tx, gomock.AssignableToTypeOf(models.Reservation{})).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
		{
			name:     "Fail reservation was accepted",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("given reservation is already accepted"),
		},
		{
			name:     "Fail reservation was cancelled",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("given reservation was cancelled"),
		},
		{
			name:     "Fail reserved reservation does not exist",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Reserved).
					Return(models.Reservation{}, errors.New("not nil"))
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("reserved reservation with given params does not exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			resRepo := mock_storage.NewMockReservation(ctrl)
			testCase.reservationRepoBehavior(resRepo, tx)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			reservationer := services.Reservationer{
				Reservation: resRepo,
				TxHelper:    txHelper,
			}

			err := reservationer.AcceptReservation(testCase.inputDto)
			assert.EqualValues(t, testCase.expectedError, err)
		})
	}
}

func TestReservationer_CancelReservation(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal)
	type reservationRepoBehavior func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                    string
		inputDto                dtos.ReservationDto
		accountRepoBehavior     accountRepoBehavior
		reservationRepoBehavior reservationRepoBehavior
		txHelperBehavior        txHelperBehavior
		expectedError           error
	}{
		{
			name:     "Success",
			inputDto: dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().IncreaseBalance(tx, accountId, value).Return(nil)
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Reserved).
					Return(models.Reservation{}, nil)
				reservationRepo.EXPECT().CreateReservation(tx, gomock.AssignableToTypeOf(models.Reservation{})).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
		{
			name:                "Fail reservation was canceled",
			inputDto:            dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("given reservation is already cancelled"),
		},
		{
			name:                "Fail reservation was accepted",
			inputDto:            dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("given reservation was accepted"),
		},
		{
			name:                "Fail reserved reservation does not exist",
			inputDto:            dtos.ReservationDto{AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservation, tx *sql.Tx) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dtos.ReservationDto{}), models.Reserved).
					Return(models.Reservation{}, errors.New("not nil"))
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("reserved reservation with given params does not exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, testCase.inputDto.AccountId, testCase.inputDto.TotalCost)

			resRepo := mock_storage.NewMockReservation(ctrl)
			testCase.reservationRepoBehavior(resRepo, tx)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			reservationer := services.Reservationer{
				Account:     accRepo,
				Reservation: resRepo,
				TxHelper:    txHelper,
			}

			err := reservationer.CancelReservation(testCase.inputDto)
			assert.EqualValues(t, testCase.expectedError, err)
		})
	}
}
