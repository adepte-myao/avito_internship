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
	Account     storage.Account
	Reservation storage.Reservation
	TxHelper    storage.SQLTransactionHelper
}

func NewReservationer(repo *storage.SQLRepository) *Reservationer {
	return &Reservationer{
		Account:     repo.Account,
		Reservation: repo.Reservation,
		TxHelper:    repo.SQLTransactionHelper,
	}
}

func (serv *Reservationer) MakeReservation(resDto dtos.ReservationDto) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	// TODO: reservation can already exist

	account, err := serv.Account.GetAccount(tx, resDto.AccountId)
	if err != nil {
		return fmt.Errorf("account with ID %d does not exist", resDto.AccountId)
	}

	if account.Balance.LessThan(resDto.TotalCost) {
		return fmt.Errorf("not enough money: account: id: %d, balance: %s; required: %s",
			account.ID, account.Balance.String(), resDto.TotalCost.String())
	}

	err = serv.Account.DecreaseBalance(tx, resDto.AccountId, resDto.TotalCost)
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
	err = serv.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Reservationer) AcceptReservation(resDto dtos.ReservationDto) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	reservation, err := serv.Reservation.GetReservation(tx, resDto, models.Accepted)
	if err == nil {
		return errors.New("given reservation is already accepted")
	}

	reservation, err = serv.Reservation.GetReservation(tx, resDto, models.Cancelled)
	if err == nil {
		return errors.New("given reservation was cancelled")
	}

	reservation, err = serv.Reservation.GetReservation(tx, resDto, models.Reserved)
	if err != nil {
		return errors.New("reserved reservation with given params does not exist")
	}

	reservation.State = models.Accepted
	reservation.RecordTime = time.Now()
	err = serv.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		return err
	}
	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Reservationer) CancelReservation(resDto dtos.ReservationDto) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	reservation, err := serv.Reservation.GetReservation(tx, resDto, models.Cancelled)
	if err == nil {
		return errors.New("given reservation is already cancelled")
	}

	reservation, err = serv.Reservation.GetReservation(tx, resDto, models.Accepted)
	if err == nil {
		return errors.New("given reservation was accepted")
	}

	reservation, err = serv.Reservation.GetReservation(tx, resDto, models.Reserved)
	if err != nil {
		return errors.New("reserved reservation with given params does not exist")
	}

	err = serv.Account.IncreaseBalance(tx, resDto.AccountId, resDto.TotalCost)
	if err != nil {
		return err
	}

	reservation.State = models.Cancelled
	reservation.RecordTime = time.Now()
	reservation.BalanceAfter = reservation.BalanceAfter.Add(resDto.TotalCost)
	err = serv.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Reservationer) GetAccountantReport() ([]models.AccountantReportElem, error) {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	report, err := serv.Reservation.GetAccountantReport(tx)
	if err != nil {
		return nil, err
	}

	serv.TxHelper.CommitTransaction(tx)

	return report, nil
}
