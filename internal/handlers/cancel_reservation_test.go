package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/handlers"
	mock_services "github.com/adepte-myao/avito_internship/internal/services/mock_service"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCancelReservationHandler(t *testing.T) {
	type reservationServBehavior func(reservationServ *mock_services.MockReservation, dto dtos.ReservationDto)

	testCases := []struct {
		name                    string
		inputBody               string
		reservationServBehavior reservationServBehavior
		expextedStatusCode      int
		expectedResponseBody    string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationServBehavior: func(reservationServ *mock_services.MockReservation, dto dtos.ReservationDto) {
				reservationServ.EXPECT().CancelReservation(gomock.AssignableToTypeOf(dtos.ReservationDto{})).
					Return(nil)
			},
			expextedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:                    "Invalid request body",
			inputBody:               `{"accountId":"1","serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationServBehavior: func(reservationServ *mock_services.MockReservation, dto dtos.ReservationDto) {},
			expextedStatusCode:      400,
			expectedResponseBody:    "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:      "Error in reservation service is not changed",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationServBehavior: func(reservationServ *mock_services.MockReservation, dto dtos.ReservationDto) {
				reservationServ.EXPECT().CancelReservation(gomock.AssignableToTypeOf(dtos.ReservationDto{})).
					Return(errors.New("bla-bla-bla"))
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"bla-bla-bla\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			resServ := mock_services.NewMockReservation(ctrl)
			testCase.reservationServBehavior(resServ, dtos.ReservationDto{
				AccountId: 1, ServiceId: 1, OrderId: 1, TotalCost: decimal.NewFromInt(100)})

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.Handler{
				Logger:      logger,
				Reservation: resServ,
			}
			router := handler.InitRoutes()

			req, err := http.NewRequest("POST", "/reservation/cancel",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()

			router.ServeHTTP(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
