package storage

import (
	"database/sql"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
)

type ReservationRepository struct{}

func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{}
}

func (repo *ReservationRepository) CreateReservation(tx *sql.Tx, reservation models.Reservation) error {
	_, err := tx.Exec(
		`INSERT INTO reserves_history (accountID, serviceID, orderID, totalCost, state, record_time, balanceAfter)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		reservation.AccountId,
		reservation.ServiceId,
		reservation.OrderId,
		reservation.TotalCost,
		reservation.State,
		reservation.RecordTime,
		reservation.BalanceAfter,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ReservationRepository) GetReservation(tx *sql.Tx, reservationDto dtos.ReservationDto) (models.Reservation, error) {
	reservation := models.Reservation{
		AccountId: reservationDto.AccountId,
		ServiceId: reservationDto.ServiceId,
		OrderId:   reservationDto.OrderId,
		TotalCost: reservationDto.TotalCost,
	}
	err := tx.QueryRow(
		`SELECT state, balanceAfter
			WHERE accountID = $1 AND serviceID = $2 AND orderID = $3 AND totalCost = $4`,
		reservationDto.AccountId,
		reservationDto.ServiceId,
		reservationDto.OrderId,
		reservationDto.TotalCost,
	).Scan(&reservation.State, &reservation.BalanceAfter)
	return reservation, err
}
