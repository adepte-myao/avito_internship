package handlers_test

import (
	"bytes"
	"database/sql"
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
	type txHelperBeginTxBehavior func(txHelper *mock_storage.MockSQLTransactionHelper)

	testCases := []struct {
		name                    string
		inputBody               string
		accRepoBehavior         accRepoBehavior
		txHelperBeginTxBehavior txHelperBeginTxBehavior
		expextedStatusCode      int
		expectedResponseBody    string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":"1"}`,
			accRepoBehavior: func(accountRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accountId int32) {
				accountRepo.EXPECT().GetAccount(tx, accountId).Return(models.Account{
					ID:      int32(1),
					Balance: decimal.NewFromInt(100),
				}, nil,
				)
			},
			txHelperBeginTxBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper) {
				txHelper.EXPECT().BeginTransaction().Return(sql.Tx{}, nil)
			},
			expextedStatusCode:   200,
			expectedResponseBody: `{"accountId":"1","balance":"100.00"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}
			accRepo := mock_storage.NewMockAccountRepo(ctrl)
			testCase.accRepoBehavior(accRepo, tx, int32(1))

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBeginTxBehavior(txHelper)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.GetBalanceHandler{
				Logger:      logger,
				AccountRepo: accRepo,
				TxHelper:    txHelper,
			}

			req, err := http.NewRequest("POST", "/get-balance",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()
			handler.Handle(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
