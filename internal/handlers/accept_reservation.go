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

type AcceptReservationHandler struct {
	Logger          *logrus.Logger
	ReservationRepo storage.ReservationRepo
	TxHelper        storage.SQLTransactionHelper
}

func NewAcceptReservationHandler(Logger *logrus.Logger, store *storage.Storage) *AcceptReservationHandler {
	return &AcceptReservationHandler{
		Logger:          Logger,
		ReservationRepo: storage.NewReservationRepository(),
		TxHelper:        storage.NewTransactionHelper(store),
	}
}

func (handler *AcceptReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Make reservation request received")

	var data dtos.ReservationDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		// TODO: find out what part of body was not decoded, maybe get more user-friendly output
		return
	}

	tx, err := handler.TxHelper.BeginTransaction()
	if err != nil {
		// TODO
		return
	}
	defer handler.TxHelper.RollbackTransaction(tx)

	reservation, err := handler.ReservationRepo.GetReservation(tx, data, models.Reserved)
	if err != nil {
		// TODO
		return
	}

	if reservation.State != models.Reserved {
		// TODO
		return
	}

	reservation.State = models.Accepted
	reservation.RecordTime = time.Now()
	err = handler.ReservationRepo.CreateReservation(tx, reservation)
	if err != nil {
		// TODO
		return
	}
	handler.TxHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
