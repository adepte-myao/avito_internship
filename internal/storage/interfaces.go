package storage

import (
	"database/sql"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AccountRepo interface {
	GetAccount(tx *sql.Tx, id int32) (models.Account, error)
	CreateAccount(tx *sql.Tx, id int32) error
	IncreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error
	DecreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error
}

type ReservationRepo interface {
	CreateReservation(tx *sql.Tx, reservation models.Reservation) error
	GetReservation(tx *sql.Tx, reservationDto dtos.ReservationDto) (models.Reservation, error)
}

type SQLTransactionHelper interface {
	BeginTransaction() (*sql.Tx, error)
	RollbackTransaction(tx *sql.Tx)
	CommitTransaction(tx *sql.Tx)
}
