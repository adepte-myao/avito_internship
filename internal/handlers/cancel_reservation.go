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

type CancelReservationHandler struct {
	logger          *logrus.Logger
	accountRepo     *storage.AccountRepository
	reservationRepo *storage.ReservationRepository
	txHelper        *storage.TransactionHelper
}

func NewCancelReservationHandler(logger *logrus.Logger, store *storage.Storage) *CancelReservationHandler {
	return &CancelReservationHandler{
		logger:          logger,
		accountRepo:     storage.NewAccountRepository(),
		reservationRepo: storage.NewReservationRepository(),
		txHelper:        storage.NewTransactionHelper(store),
	}
}

func (handler *CancelReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
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

	reservation, err := handler.reservationRepo.GetReservation(tx, data, models.Reserved)
	if err != nil {
		// TODO
		return
	}

	reservation.State = models.Cancelled
	reservation.RecordTime = time.Now()
	err = handler.reservationRepo.CreateReservation(tx, reservation)
	if err != nil {
		// TODO
		return
	}

	handler.txHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
