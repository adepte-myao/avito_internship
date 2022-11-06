package storage

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

type ReservationRepository struct{}

func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{}
}

func (repo *ReservationRepository) CreateReservation(tx *sql.Tx, reservation models.Reservation) error {
	_, err := tx.Exec(
		`INSERT INTO reserves_history (account_id, service_id, order_id, total_cost, state, record_time, balance_after)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		reservation.AccountId,
		reservation.ServiceId,
		reservation.OrderId,
		strings.Replace(reservation.TotalCost.String(), ".", ",", -1),
		reservation.State.String(),
		reservation.RecordTime,
		strings.Replace(reservation.BalanceAfter.String(), ".", ",", -1),
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ReservationRepository) GetReservation(tx *sql.Tx, reservationDto dtos.ReservationDto, state models.ReserveState) (models.Reservation, error) {
	reservation := models.Reservation{
		AccountId: reservationDto.AccountId,
		ServiceId: reservationDto.ServiceId,
		OrderId:   reservationDto.OrderId,
		TotalCost: reservationDto.TotalCost,
	}
	var balance string
	err := tx.QueryRow(
		`SELECT balance_after FROM reserves_history
			WHERE account_id = $1 AND service_id = $2 AND order_id = $3 AND total_cost = $4 AND state = $5`,
		reservationDto.AccountId,
		reservationDto.ServiceId,
		reservationDto.OrderId,
		strings.Replace(reservationDto.TotalCost.String(), ".", ",", -1),
		state.String(),
	).Scan(&balance)

	if err != nil {
		return models.Reservation{}, err
	}

	// balance format: "123.45 P"
	regBalance := regexp.MustCompile(`[^0-9,]`)
	balance = regBalance.ReplaceAllString(balance, "")

	balance = strings.Replace(balance, ",", ".", 1)

	reservation.BalanceAfter, err = decimal.NewFromString(balance)
	if err != nil {
		return models.Reservation{}, err
	}

	return reservation, err
}
