package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/handlers"
	mock_services "github.com/adepte-myao/avito_internship/internal/services/mock_service"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawAccountHandler(t *testing.T) {
	type accountServBehavior func(accServ *mock_services.MockAccount, accId int32, value decimal.Decimal)

	testCases := []struct {
		name                 string
		inputBody            string
		accountServBehavior  accountServBehavior
		expextedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":1,"value":"100.01"}`,
			accountServBehavior: func(accServ *mock_services.MockAccount, accId int32, value decimal.Decimal) {
				accServ.EXPECT().Withdraw(accId, value).Return(nil)
			},
			expextedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:                 "Invalid request body",
			inputBody:            `{"accountId":"1","value":"100.01"}`,
			accountServBehavior:  func(accServ *mock_services.MockAccount, accId int32, value decimal.Decimal) {},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:      "Error from reservation service is not changed",
			inputBody: `{"accountId":1,"value":"100.01"}`,
			accountServBehavior: func(accServ *mock_services.MockAccount, accId int32, value decimal.Decimal) {
				accServ.EXPECT().Withdraw(accId, value).Return(errors.New("bla-bla-bla"))
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"bla-bla-bla\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			accServ := mock_services.NewMockAccount(ctrl)
			testCase.accountServBehavior(accServ, int32(1), decimal.NewFromFloat(100.01))

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.Handler{
				Logger:  logger,
				Account: accServ,
			}
			router := handler.InitRoutes()

			req, err := http.NewRequest("POST", "/balance/withdraw",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()

			router.ServeHTTP(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
