package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type MakeReservationHandler struct {
	logger          *logrus.Logger
	accountRepo     *storage.AccountRepository
	reservationRepo *storage.ReservationRepository
	txHelper        *storage.TransactionHelper
}

func NewMakeReservationHandler(logger *logrus.Logger, store *storage.Storage) *MakeReservationHandler {
	return &MakeReservationHandler{
		logger:          logger,
		accountRepo:     storage.NewAccountRepository(),
		reservationRepo: storage.NewReservationRepository(),
		txHelper:        storage.NewTransactionHelper(store),
	}
}

func (handler *MakeReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Make reservation request received")

	var data dtos.ReservationDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.logger.Error("cannot decode request body: ", err.Error())
		// TODO: find out what part of body was not decoded, maybe get more user-friendly output
		return
	}

	tx, err := handler.txHelper.BeginTransaction()
	if err != nil {
		// TODO
		return
	}
	defer handler.txHelper.RollbackTransaction(tx)

	account, err := handler.accountRepo.GetAccount(tx, data.AccountId)
	if err != nil {
		// TODO
		return
	}

	if account.Balance.LessThan(data.TotalCost) {
		// TODO
		return
	}

	err = handler.accountRepo.DecreaseBalance(tx, data.AccountId, data.TotalCost)
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
	err = handler.reservationRepo.CreateReservation(tx, reservation)
	if err != nil {
		// TODO
		return
	}

	handler.txHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
