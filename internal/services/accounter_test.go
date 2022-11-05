package services_test

import (
	"database/sql"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/adepte-myao/avito_internship/internal/storage/mock_storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccounter_GetBalance(t *testing.T) {
	type accountRepoBehavior func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, dto dtos.GetBalanceDto)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                string
		inputDto            dtos.GetBalanceDto
		accountRepoBehavior accountRepoBehavior
		txHelperBehavior    txHelperBehavior
		expextedBalance     decimal.Decimal
		expectedError       error
	}{
		{
			name:     "Success",
			inputDto: dtos.GetBalanceDto{AccountId: 1},
			accountRepoBehavior: func(accountRepo *mock_storage.MockAccount, tx *sql.Tx, dto dtos.GetBalanceDto) {
				accountRepo.EXPECT().GetAccount(tx, dto.AccountId).Return(
					models.Account{ID: dto.AccountId, Balance: decimal.NewFromInt(100)}, nil,
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}
			dto := dtos.GetBalanceDto{AccountId: 1}

			accRepo := mock_storage.NewMockAccount(ctrl)
			testCase.accountRepoBehavior(accRepo, tx, dto)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			accounter := services.Accounter{
				Account:  accRepo,
				TxHelper: txHelper,
			}

			account, err := accounter.GetBalance(dto)
			assert.EqualValues(t, testCase.expectedError, err)
			if err == nil {
				assert.Equal(t, account.ID, dto.AccountId)
				isBalanceCorrect := account.Balance.Equal(testCase.expextedBalance)
				assert.True(t, isBalanceCorrect)
			}
		})
	}
}
