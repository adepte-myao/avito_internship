package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/models"
	mock_services "github.com/adepte-myao/avito_internship/internal/services/mock_service"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountantReportHandler(t *testing.T) {
	type reservationServBehavior func(resServ *mock_services.MockReservation)

	testCases := []struct {
		name                    string
		inputBody               string
		reservationServBehavior reservationServBehavior
		expextedStatusCode      int
		expectedResponseBody    string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":1}`,
			reservationServBehavior: func(resServ *mock_services.MockReservation) {
				report := make([]models.AccountantReportElem, 0)
				elem := models.AccountantReportElem{
					ServiceName:   "bla-bla",
					TotalReceived: decimal.NewFromFloat(100.01),
				}
				report = append(report, elem)
				resServ.EXPECT().GetAccountantReport(11, 2022).Return(report, nil)
			},
			expextedStatusCode:   200,
			expectedResponseBody: "[{\"serviceName\":\"bla-bla\",\"totalReceived\":\"100.01\"}]\n",
		},
		{
			name:      "Error from reservation service is not changed",
			inputBody: `{"accountId":1}`,
			reservationServBehavior: func(resServ *mock_services.MockReservation) {
				resServ.EXPECT().GetAccountantReport(11, 2022).Return(nil, errors.New("bla-bla-bla"))
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"bla-bla-bla\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			reServ := mock_services.NewMockReservation(ctrl)
			testCase.reservationServBehavior(reServ)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.Handler{
				Logger:      logger,
				Reservation: reServ,
			}
			router := handler.InitRoutes()

			req, err := http.NewRequest("GET", "/accountant-report",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()

			router.ServeHTTP(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
