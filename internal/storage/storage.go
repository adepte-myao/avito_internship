package storage

import (
	"database/sql"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=storage.go -destination=mock_storage/mock.go

type Account interface {
	GetAccount(tx *sql.Tx, id int32) (models.Account, error)
	CreateAccount(tx *sql.Tx, id int32) error
	IncreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error
	DecreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error
}

type Reservation interface {
	CreateReservation(tx *sql.Tx, reservation models.Reservation) error
	GetReservation(tx *sql.Tx, reservationDto dtos.ReservationDto, state models.ReserveState) (models.Reservation, error)
}

type Transfer interface {
	RecordExternalTransfer(tx *sql.Tx, accId int32, ttype models.TransferType, amount decimal.Decimal) error
	RecordInternalTransfer(tx *sql.Tx, senderId int32, recId int32, amount decimal.Decimal) error
}

type SQLTransactionHelper interface {
	BeginTransaction() (*sql.Tx, error)
	RollbackTransaction(tx *sql.Tx)
	CommitTransaction(tx *sql.Tx)
}

type SQLRepository struct {
	Account
	Reservation
	Transfer
	SQLTransactionHelper
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{
		Account:              NewAccountRepository(),
		Reservation:          NewReservationRepository(),
		Transfer:             NewTransferer(),
		SQLTransactionHelper: NewTransactionHelper(db),
	}
}
