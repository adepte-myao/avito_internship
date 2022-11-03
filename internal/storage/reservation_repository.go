package storage

import (
	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
)

type ReservationRepository struct {
	storage *Storage
}

func NewReservationRepository(storage *Storage) *ReservationRepository {
	return &ReservationRepository{
		storage: storage,
	}
}

func (repo *ReservationRepository) CreateReservation(reservation models.Reservation) error {
	tx, err := repo.storage.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`UPDATE accounts 
			SET balance = balance - $1 
			WHERE id = $2`,
		reservation.TotalCost,
		reservation.AccountId,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
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

	tx.Commit()
	return nil
}

func (repo *ReservationRepository) GetReservation(reservationDto dtos.ReservationDto) (models.Reservation, error) {
	reservation := models.Reservation{
		AccountId: reservationDto.AccountId,
		ServiceId: reservationDto.ServiceId,
		OrderId:   reservationDto.OrderId,
		TotalCost: reservationDto.TotalCost,
	}
	err := repo.storage.db.QueryRow(
		`SELECT state, balanceAfter
			WHERE accountID = $1 AND serviceID = $2 AND orderID = $3 AND totalCost = $4`,
		reservationDto.AccountId,
		reservationDto.ServiceId,
		reservationDto.OrderId,
		reservationDto.TotalCost,
	).Scan(&reservation.State, &reservation.BalanceAfter)
	return reservation, err
}
