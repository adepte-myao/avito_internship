package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type MakeReservationHandler struct {
	Logger          *logrus.Logger
	AccountRepo     storage.AccountRepo
	ReservationRepo storage.ReservationRepo
	TxHelper        storage.SQLTransactionHelper
}

func NewMakeReservationHandler(Logger *logrus.Logger, store *storage.Storage) *MakeReservationHandler {
	return &MakeReservationHandler{
		Logger:          Logger,
		AccountRepo:     storage.NewAccountRepository(),
		ReservationRepo: storage.NewReservationRepository(),
		TxHelper:        storage.NewTransactionHelper(store),
	}
}

func (handler *MakeReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Make reservation request received")

	var data dtos.ReservationDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "invalid request body",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	tx, err := handler.TxHelper.BeginTransaction()
	if err != nil {
		// Should not be there
		return
	}
	defer handler.TxHelper.RollbackTransaction(tx)

	account, err := handler.AccountRepo.GetAccount(tx, data.AccountId)
	if err != nil {
		// TODO
		return
	}

	if account.Balance.LessThan(data.TotalCost) {
		// TODO
		return
	}

	err = handler.AccountRepo.DecreaseBalance(tx, data.AccountId, data.TotalCost)
	if err != nil {
		// TODO
		return
	}

	reservation := models.Reservation{
		AccountId:    data.AccountId,
		ServiceId:    data.ServiceId,
		OrderId:      data.OrderId,
		TotalCost:    data.TotalCost,
		State:        models.Reserved,
		RecordTime:   time.Now(),
		BalanceAfter: account.Balance.Sub(data.TotalCost),
	}
	err = handler.ReservationRepo.CreateReservation(tx, reservation)
	if err != nil {
		// TODO
		return
	}

	handler.TxHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
