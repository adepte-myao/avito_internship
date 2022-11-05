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

func TestGetBalanceHandler(t *testing.T) {
	type accRepoBehavior func(accountRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accountId int32)
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
			name:      "Success without fractional part",
			inputBody: `{"accountId":1}`,
			accRepoBehavior: func(accountRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(models.Account{
					ID:      int32(1),
					Balance: decimal.NewFromInt(100),
				}, nil,
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   200,
			expectedResponseBody: "{\"accountId\":1,\"balance\":\"100\"}\n",
		},
		{
			name:      "Success with fractional part",
			inputBody: `{"accountId":1}`,
			accRepoBehavior: func(accountRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(models.Account{
					ID:      int32(1),
					Balance: decimal.NewFromFloat(100.01),
				}, nil,
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().CommitTransaction(tx).Return()
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   200,
			expectedResponseBody: "{\"accountId\":1,\"balance\":\"100.01\"}\n",
		},
		{
			name:                 "Invalid request body",
			inputBody:            `{"accountId":"1"}`,
			accRepoBehavior:      func(accountRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accountId int32) {},
			txHelperBehavior:     func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:      "Account does not exist",
			inputBody: `{"accountId":1}`,
			accRepoBehavior: func(accountRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(models.Account{},
					errors.New("not a real error, but it's not nil"),
				)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"account with given ID does not exist\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}
			accRepo := mock_storage.NewMockAccountRepo(ctrl)
			testCase.accRepoBehavior(accRepo, tx, int32(1))

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.GetBalanceHandler{
				Logger:      logger,
			}

			req, err := http.NewRequest("GET", "/get-balance",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()
			handler.Handle(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
