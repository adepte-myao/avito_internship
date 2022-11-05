package handlers_test

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage/mock_storage"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAcceptReservationHandler(t *testing.T) {
	type reservationRepoBehavior func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
		dto dtos.ReservationDto, reservation models.Reservation)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                    string
		inputBody               string
		reservationRepoBehavior reservationRepoBehavior
		txHelperBehavior        txHelperBehavior
		expextedStatusCode      int
		expectedResponseBody    string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil error"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil error"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Reserved).
					Return(models.Reservation{State: models.Reserved}, nil)
				reservationRepo.EXPECT().CreateReservation(tx, gomock.AssignableToTypeOf(reservation)).Return(nil)
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
			name:      "Invalid request body",
			inputBody: `{"accountId":"1","serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
			},
			txHelperBehavior:     func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:      "Reservation does not exist",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil error"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Cancelled).
					Return(models.Reservation{}, errors.New("not nil error"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Reserved).
					Return(models.Reservation{}, errors.New("not nil error"))
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"reserved reservation with given params does not exist\"}\n",
		},
		{
			name:      "Reservation already accepted",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Accepted).
					Return(models.Reservation{}, nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"given reservation is already accepted\"}\n",
		},
		{
			name:      "Reservation was cancelled",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Accepted).
					Return(models.Reservation{}, errors.New("not nil"))
				reservationRepo.EXPECT().GetReservation(tx, gomock.AssignableToTypeOf(dto), models.Cancelled).
					Return(models.Reservation{}, nil)
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"given reservation was cancelled\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := &sql.Tx{}
			// Dont't need specific values here because in calls we just check types
			dto := dtos.ReservationDto{}
			reservation := models.Reservation{}

			reservationRepo := mock_storage.NewMockReservationRepo(ctrl)
			testCase.reservationRepoBehavior(reservationRepo, tx, dto, reservation)

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.AcceptReservationHandler{
				Logger:          logger,
			}

			req, err := http.NewRequest("POST", "/accept-reservation",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()
			handler.Handle(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}
