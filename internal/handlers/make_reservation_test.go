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
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMakeReservationHandler(t *testing.T) {
	type accRepoBehavior func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal)
	type reservationRepoBehavior func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
		dto dtos.ReservationDto, reservation models.Reservation)
	type txHelperBehavior func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx)

	testCases := []struct {
		name                    string
		inputBody               string
		accRepoBehavior         accRepoBehavior
		reservationRepoBehavior reservationRepoBehavior
		txHelperBehavior        txHelperBehavior
		expextedStatusCode      int
		expectedResponseBody    string
	}{
		{
			name:      "Success",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			accRepoBehavior: func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal) {
				accRepo.EXPECT().GetAccount(tx, accId).Return(models.Account{ID: 1, Balance: decimal.NewFromInt(200)}, nil)
				accRepo.EXPECT().DecreaseBalance(tx, accId, gomock.AssignableToTypeOf(decimal.NewFromInt(1))).Return(nil)
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
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
			name:      "Account does not exist",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			accRepoBehavior: func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal) {
				accRepo.EXPECT().GetAccount(tx, accId).Return(models.Account{}, errors.New("not nil"))
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"account does not exist\"}\n",
		},
		{
			name:      "Not enough money",
			inputBody: `{"accountId":1,"serviceId":1,"orderId":1,"totalCost":"100.00"}`,
			accRepoBehavior: func(accRepo *mock_storage.MockAccountRepo, tx *sql.Tx, accId int32, totalCost decimal.Decimal) {
				accRepo.EXPECT().GetAccount(tx, accId).Return(models.Account{ID: 1, Balance: decimal.NewFromInt(50)}, nil)
			},
			reservationRepoBehavior: func(reservationRepo *mock_storage.MockReservationRepo, tx *sql.Tx,
				dto dtos.ReservationDto, reservation models.Reservation) {
			},
			txHelperBehavior: func(txHelper *mock_storage.MockSQLTransactionHelper, tx *sql.Tx) {
				txHelper.EXPECT().BeginTransaction().Return(&sql.Tx{}, nil)
				txHelper.EXPECT().RollbackTransaction(tx).Return()
			},
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"not enough money\"}\n",
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

			accRepo := mock_storage.NewMockAccountRepo(ctrl)
			testCase.accRepoBehavior(accRepo, tx, 1, decimal.NewFromInt(100))

			txHelper := mock_storage.NewMockSQLTransactionHelper(ctrl)
			testCase.txHelperBehavior(txHelper, tx)

			logger := logrus.New()
			logger.Level = logrus.FatalLevel

			handler := handlers.MakeReservationHandler{
				Logger:          logger,
				AccountRepo:     accRepo,
				ReservationRepo: reservationRepo,
				TxHelper:        txHelper,
			}

			req, err := http.NewRequest("POST", "/make-reservation",
				bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			rw := httptest.NewRecorder()
			handler.Handle(rw, req)

			assert.Equal(t, testCase.expextedStatusCode, rw.Code)
			assert.Equal(t, testCase.expectedResponseBody, rw.Body.String())
		})
	}
}