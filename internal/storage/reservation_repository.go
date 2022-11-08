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

func (repo *ReservationRepository) GetAccountantReport(tx *sql.Tx, month int, year int) ([]models.AccountantReportElem, error) {
	rows, err := tx.Query(
		`SELECT s.name, total
		FROM (SELECT service_id, sum(total_cost) AS "total"
			  FROM reserves_history
			  WHERE state = 'accepted'::reserve_state AND
			        date_part('month', record_time) = $1 AND
			        date_part('year', record_time) = $2
			  GROUP BY service_id) AS t
		JOIN services s ON service_id = s.id`,
		month, year)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	report := make([]models.AccountantReportElem, 0)
	for rows.Next() {
		elem := models.AccountantReportElem{}
		var total string
		err = rows.Scan(&elem.ServiceName, &total)
		if err != nil {
			return nil, err
		}

		// total format: "123.45 P"
		regTotal := regexp.MustCompile(`[^0-9,]`)
		total = regTotal.ReplaceAllString(total, "")
		total = strings.Replace(total, ",", ".", 1)
		elem.TotalReceived, err = decimal.NewFromString(total)
		if err != nil {
			return nil, err
		}

		report = append(report, elem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}
