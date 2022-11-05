package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
)

type Reservationer struct {
	repo *storage.SQLRepository
}

func NewReservationer(repo *storage.SQLRepository) *Reservationer {
	return &Reservationer{
		repo: repo,
	}
}

func (serv *Reservationer) MakeReservation(resDto dtos.ReservationDto) error {
	tx, err := serv.repo.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.repo.SQLTransactionHelper.RollbackTransaction(tx)

	account, err := serv.repo.Account.GetAccount(tx, resDto.AccountId)
	if err != nil {
		return fmt.Errorf("account with ID %d does not exist", resDto.AccountId)
	}

	if account.Balance.LessThan(resDto.TotalCost) {
		return fmt.Errorf("not enough money: account: id: %d, balance: %s; required: %s",
			account.ID, account.Balance.String(), resDto.TotalCost.String())
	}

	err = serv.repo.Account.DecreaseBalance(tx, resDto.AccountId, resDto.TotalCost)
	if err != nil {
		return err
	}

	reservation := models.Reservation{
		AccountId:    resDto.AccountId,
		ServiceId:    resDto.ServiceId,
		OrderId:      resDto.OrderId,
		TotalCost:    resDto.TotalCost,
		State:        models.Reserved,
		RecordTime:   time.Now(),
		BalanceAfter: account.Balance.Sub(resDto.TotalCost),
	}
	err = serv.repo.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		return err
	}

	serv.repo.SQLTransactionHelper.CommitTransaction(tx)

	return nil
}

func (serv *Reservationer) AcceptReservation(resDto dtos.ReservationDto) error {
	tx, err := serv.repo.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.repo.SQLTransactionHelper.RollbackTransaction(tx)

	reservation, err := serv.repo.Reservation.GetReservation(tx, resDto, models.Accepted)
	if err == nil {
		return errors.New("given reservation is already accepted")
	}

	reservation, err = serv.repo.Reservation.GetReservation(tx, resDto, models.Cancelled)
	if err == nil {
		return errors.New("given reservation was cancelled")
	}

	reservation, err = serv.repo.Reservation.GetReservation(tx, resDto, models.Reserved)
	if err != nil {
		return errors.New("reserved reservation with given params does not exist")
	}

	reservation.State = models.Accepted
	reservation.RecordTime = time.Now()
	err = serv.repo.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		return err
	}
	serv.repo.SQLTransactionHelper.CommitTransaction(tx)

	return nil
}

func (serv *Reservationer) CancelReservation(resDto dtos.ReservationDto) error {
	tx, err := serv.repo.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.repo.SQLTransactionHelper.RollbackTransaction(tx)

	reservation, err := serv.repo.Reservation.GetReservation(tx, resDto, models.Cancelled)
	if err == nil {
		return errors.New("given reservation is already cancelled")
	}

	reservation, err = serv.repo.Reservation.GetReservation(tx, resDto, models.Accepted)
	if err == nil {
		return errors.New("given reservation was accepted")
	}

	reservation, err = serv.repo.Reservation.GetReservation(tx, resDto, models.Reserved)
	if err != nil {
		return errors.New("reserved reservation with given params does not exist")
	}

	err = serv.repo.Account.IncreaseBalance(tx, reservation.AccountId, reservation.TotalCost)
	if err != nil {
		return err
	}

	reservation.State = models.Cancelled
	reservation.RecordTime = time.Now()
	err = serv.repo.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		return err
	}

	serv.repo.SQLTransactionHelper.CommitTransaction(tx)

	return nil
}
