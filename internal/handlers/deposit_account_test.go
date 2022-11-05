package handlers_test

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage/mock_storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDepositAccountHandler(t *testing.T) {
	type accRepoBehavior func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                 string
		inputBody            string
		accRepoBehavior      accRepoBehavior
		txHelperBehavior     txHelperBehavior
		expextedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Success account exists",
			inputBody: `{"accountId":1,"value":"100.00"}`,
			accRepoBehavior: func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal) {
				accRepo.EXPECT().GetAccount(tx, accId).Return(models.Account{ID: 1, Balance: decimal.NewFromInt(200)}, nil)
				accRepo.EXPECT().IncreaseBalance(tx, accId, gomock.AssignableToTypeOf(decimal.NewFromInt(1))).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:      "Success account does not exist",
			inputBody: `{"accountId":1,"value":"100.00"}`,
			accRepoBehavior: func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal) {
				accRepo.EXPECT().GetAccount(tx, accId).Return(models.Account{}, errors.New("not nil"))
				accRepo.EXPECT().CreateAccount(tx, accId).Return(nil)
				accRepo.EXPECT().IncreaseBalance(tx, accId, gomock.AssignableToTypeOf(decimal.NewFromInt(1))).Return(nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:                 "Invalid request body",
			inputBody:            `{"accountId":"1","value":"100.00"}`,
			accRepoBehavior:      func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal) {},
			txHelperBehavior:     func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}

			accRepo := mock_storage.NewMockAccountRepo(ctrl)
			testCase.accRepoBehavior(accRepo, tx, 1, decimal.NewFromInt(100))

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.DepositAccountHandler{
				Logger:      logger,
			}

			req, err := http.NewRequest("POST", "/deposit-account",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()
			handler.Handle(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
