package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/models"
	mock_services "github.com/adepte-myao/avito_internship/internal/services/mock_service"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountStatementHandler(t *testing.T) {
	type accountServBehavior func(accServ *mock_services.MockAccount)

	testCases := []struct {
		name                 string
		inputBody            string
		accountServBehavior  accountServBehavior
		expextedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":1,"page":1,"pageSize":1,"firstSortCriteria":"","secondSortCriteria":""}`,
			accountServBehavior: func(accServ *mock_services.MockAccount) {
				statements := make([]models.StatementElem, 0)
				statement := models.StatementElem{
					RecordTime: time.Date(2022, 12, 12, 11, 10, 9, 8, time.UTC),
					TransferType: "deposit",
					Amount: decimal.NewFromInt(100),
					Description: "description example",
				}
				statements = append(statements, statement)
				dto := dtos.GetAccountStatementDto{
					AccountId: 1,
					Page: 1,
					PageSize: 1,
					FirstSortCriteria: "",
					SecondSortCriteria: "",
				}
				accServ.EXPECT().GetStatement(dto).Return(statements, nil)
			},
			expextedStatusCode:   200,
			expectedResponseBody: "[{\"recordTime\":\"2022-12-12T11:10:09.000000008Z\",\"transferType\":\"deposit\",\"amount\":\"100\",\"description\":\"description example\"}]\n",
		},
		{
			name:                 "Invalid request body",
			inputBody:            `{"accountId":"1"}`,
			accountServBehavior:  func(accServ *mock_services.MockAccount) {},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:      "Error from reservation service is not changed",
			inputBody: `{"accountId":1,"page":1,"pageSize":1,"firstSortCriteria":"","secondSortCriteria":""}`,
			accountServBehavior: func(accServ *mock_services.MockAccount) {
				dto := dtos.GetAccountStatementDto{}
				accServ.EXPECT().GetStatement(gomock.AssignableToTypeOf(dto)).Return(nil, errors.New("bla-bla-bla"))
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"bla-bla-bla\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			accServ := mock_services.NewMockAccount(ctrl)
			testCase.accountServBehavior(accServ)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.Handler{
				Logger:  logger,
				Account: accServ,
			}
			router := handler.InitRoutes()

			req, err := http.NewRequest("GET", "/balance/statement",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()

			router.ServeHTTP(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
