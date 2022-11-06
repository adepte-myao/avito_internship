package services_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/adepte-myao/avito_internship/internal/storage/mock_storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccounter_GetBalance(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                string
		inputId             int32
		accountRepoBehavior accountRepoBehavior
		txHelperBehavior    txHelperBehavior
		expextedBalance     decimal.Decimal
		expectedError       error
	}{
		{
			name:    "Success without fractional part",
			inputId: 1,
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromInt(100)}, nil,
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedBalance: decimal.NewFromInt(100),
			expectedError:   nil,
		},
		{
			name:    "Success with fractional part",
			inputId: 1,
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromFloat(100.02)}, nil,
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedBalance: decimal.NewFromFloat(100.02),
			expectedError:   nil,
		},
		{
			name:    "Account does not exist",
			inputId: 1,
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{}, errors.New("not nil"),
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedBalance: decimal.NewFromInt(0),
			expectedError:   errors.New("Account with given ID does not exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, testCase.inputId)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			accounter := services.Accounter{
				Account:  accRepo,
				TxHelper: txHelper,
			}

			account, err := accounter.GetBalance(testCase.inputId)
			assert.EqualValues(t, testCase.expectedError, err)
			if err == nil {
				assert.Equal(t, account.ID, testCase.inputId)
				isBalanceCorrect := account.Balance.Equal(testCase.expextedBalance)
				assert.True(t, isBalanceCorrect)
			}
		})
	}
}

func TestAccounter_Deposit(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                string
		inputAccountId      int32
		inputValue          decimal.Decimal
		accountRepoBehavior accountRepoBehavior
		txHelperBehavior    txHelperBehavior
		expectedError       error
	}{
		{
			name:           "Success account does not exist",
			inputAccountId: 1,
			inputValue:     decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{}, errors.New("not nil"),
				)
				accountRepo.EXPECT().CreateAccount(tx, accountId).Return(nil)
				accountRepo.EXPECT().IncreaseBalance(tx, accountId, value).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
		{
			name:           "Success account exists",
			inputAccountId: 1,
			inputValue:     decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromInt(0)}, nil,
				)
				accountRepo.EXPECT().IncreaseBalance(tx, accountId, value).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, testCase.inputAccountId, testCase.inputValue)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			accounter := services.Accounter{
				Account:  accRepo,
				TxHelper: txHelper,
			}

			err := accounter.Deposit(testCase.inputAccountId, testCase.inputValue)
			assert.EqualValues(t, testCase.expectedError, err)
		})
	}
}

func TestAccounter_Withdraw(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                string
		inputAccountId      int32
		inputValue          decimal.Decimal
		accountRepoBehavior accountRepoBehavior
		txHelperBehavior    txHelperBehavior
		expectedError       error
	}{
		{
			name:           "Success",
			inputAccountId: 1,
			inputValue:     decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromInt(100)}, nil,
				)
				accountRepo.EXPECT().DecreaseBalance(tx, accountId, value).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
		{
			name:           "Fail account does not exist",
			inputAccountId: 1,
			inputValue:     decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{}, errors.New("not nil"),
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("account does not exist"),
		},
		{
			name:           "Fail not enough money",
			inputAccountId: 1,
			inputValue:     decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, accountId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(
					models.Account{ID: accountId, Balance: decimal.NewFromInt(99)}, nil,
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("not enough money"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, testCase.inputAccountId, testCase.inputValue)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			accounter := services.Accounter{
				Account:  accRepo,
				TxHelper: txHelper,
			}

			err := accounter.Withdraw(testCase.inputAccountId, testCase.inputValue)
			assert.EqualValues(t, testCase.expectedError, err)
		})
	}
}

func TestAccounter_InternalTransfer(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, senderId int32, recId int32, value decimal.Decimal)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                string
		senderId            int32
		recId               int32
		value               decimal.Decimal
		accountRepoBehavior accountRepoBehavior
		txHelperBehavior    txHelperBehavior
		expectedError       error
	}{
		{
			name:     "Success",
			senderId: 1,
			recId:    2,
			value:    decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, senderId int32, recId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, senderId).Return(
					models.Account{ID: senderId, Balance: decimal.NewFromInt(100)}, nil,
				)
				accountRepo.EXPECT().DecreaseBalance(tx, senderId, value).Return(nil)
				accountRepo.EXPECT().IncreaseBalance(tx, recId, value).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: nil,
		},
		{
			name:     "Sender account does not exist",
			senderId: 1,
			recId:    2,
			value:    decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, senderId int32, recId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, senderId).Return(
					models.Account{}, errors.New("not nil"),
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("sender account does not exist"),
		},
		{
			name:     "Sender account does not have enough money",
			senderId: 1,
			recId:    2,
			value:    decimal.NewFromInt(100),
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, senderId int32, recId int32, value decimal.Decimal) {
				accountRepo.EXPECT().GetAccount(tx, senderId).Return(
					models.Account{ID: senderId, Balance: decimal.NewFromInt(99)}, nil,
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expectedError: errors.New("not enough money"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, testCase.senderId, testCase.recId, testCase.value)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			accounter := services.Accounter{
				Account:  accRepo,
				TxHelper: txHelper,
			}

			err := accounter.InternalTransfer(testCase.senderId, testCase.recId, testCase.value)
			assert.EqualValues(t, testCase.expectedError, err)
		})
	}
}
