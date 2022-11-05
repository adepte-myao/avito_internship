package storage_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func getReservationRepoTxCommitFunc(t *testing.T) (*storage.ReservationRepository, *sql.Tx, func(*sql.Tx)) {
	storageStruct, closeFunc := storage.TestStore(t, databaseURL)

	repo := storage.NewReservationRepository()
	txHelper := storage.NewTransactionHelper(storageStruct)

	tx, err := txHelper.BeginTransaction()
	assert.NoError(t, err)

	accRepo := storage.NewAccountRepository()
	var accountId int32 = 15
	err = accRepo.CreateAccount(tx, accountId)
	assert.NoError(t, err)

	return repo, tx, func(tx *sql.Tx) {
		txHelper.CommitTransaction(tx)
		closeFunc()
	}
}

func createReservation(repo *storage.ReservationRepository, tx *sql.Tx) {
	var accountId int32 = 15
	var serviceId int32 = 1
	var orderId int32 = 1
	reservation := models.Reservation{
		AccountId:    accountId,
		ServiceId:    serviceId,
		OrderId:      orderId,
		TotalCost:    decimal.NewFromInt(100),
		State:        models.Reserved,
		RecordTime:   time.Now(),
		BalanceAfter: decimal.NewFromInt(0),
	}
	repo.CreateReservation(tx, reservation)
}

func TestReservationRepository_CreateReservation_Success(t *testing.T) {
	repo, tx, commitFunc := getReservationRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 15
	var serviceId int32 = 1
	var orderId int32 = 1
	reservation := models.Reservation{
		AccountId:    accountId,
		ServiceId:    serviceId,
		OrderId:      orderId,
		TotalCost:    decimal.NewFromInt(100),
		State:        models.Reserved,
		RecordTime:   time.Now(),
		BalanceAfter: decimal.NewFromInt(0),
	}
	err := repo.CreateReservation(tx, reservation)

	assert.NoError(t, err)
}

func TestReservationRepository_GetReservation_Success(t *testing.T) {
	repo, tx, commitFunc := getReservationRepoTxCommitFunc(t)
	defer commitFunc(tx)

	createReservation(repo, tx)

	var accountId int32 = 15
	var serviceId int32 = 1
	var orderId int32 = 1
	totalCost := decimal.NewFromInt(100)
	reservationDto := dtos.ReservationDto{
		AccountId: accountId,
		ServiceId: serviceId,
		OrderId:   orderId,
		TotalCost: totalCost,
	}
	reservation, err := repo.GetReservation(tx, reservationDto, models.Reserved)

	assert.NoError(t, err)
	assert.NotNil(t, reservation)

	assert.Equal(t, models.Reserved, reservation.State)
	balanceAfterEqual := reservation.BalanceAfter.Equal(decimal.NewFromInt(0))
	assert.True(t, balanceAfterEqual)
}

func TestReservationRepository_GetReservation_FailsOnWrongParam(t *testing.T) {
	repo, tx, commitFunc := getReservationRepoTxCommitFunc(t)
	defer commitFunc(tx)

	createReservation(repo, tx)

	var accountId int32 = 15
	var serviceId int32 = 1
	var orderId int32 = 1
	totalCost := decimal.NewFromInt(100)
	reservationDto := dtos.ReservationDto{
		AccountId: accountId,
		ServiceId: serviceId,
		OrderId:   orderId,
		TotalCost: totalCost,
	}

	cases := []struct {
		name        string
		reservation func() dtos.ReservationDto
	}{
		{
			name: "wring accountId",
			reservation: func() dtos.ReservationDto {
				dto := reservationDto
				dto.AccountId = 16
				return dto
			},
		},
		{
			name: "wrong serviceId",
			reservation: func() dtos.ReservationDto {
				dto := reservationDto
				dto.ServiceId = 2
				return dto
			},
		},
		{
			name: "wrong orderId",
			reservation: func() dtos.ReservationDto {
				dto := reservationDto
				dto.OrderId = 2
				return dto
			},
		},
		{
			name: "wrong totalCost",
			reservation: func() dtos.ReservationDto {
				dto := reservationDto
				dto.TotalCost = decimal.NewFromInt(12345)
				return dto
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.GetReservation(tx, testCase.reservation(), models.Reserved)
			assert.Error(t, err)
		})
	}

}
