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

type CancelReservationHandler struct {
	Logger          *logrus.Logger
	ReservationRepo storage.ReservationRepo
	TxHelper        storage.SQLTransactionHelper
}

func NewCancelReservationHandler(Logger *logrus.Logger, store *storage.Storage) *CancelReservationHandler {
	return &CancelReservationHandler{
		Logger:          Logger,
		ReservationRepo: storage.NewReservationRepository(),
		TxHelper:        storage.NewTransactionHelper(store),
	}
}

func (handler *CancelReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
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
		// Shouldn't be there
		return
	}
	defer handler.TxHelper.RollbackTransaction(tx)

	reservation, err := handler.ReservationRepo.GetReservation(tx, data, models.Reserved)
	if err != nil {
		handler.Logger.Errorf("no reserved reservations with params: accountID: %d, serviceID: %d, orderID: %d, totalCost: %s exist",
			data.AccountId, data.ServiceId, data.OrderId, data.TotalCost.String())

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "reserved reservation with given params does not exist",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	reservation, err = handler.ReservationRepo.GetReservation(tx, data, models.Cancelled)
	if err == nil {
		handler.Logger.Errorf("already cancelled reservation with params: accountID: %d, serviceID: %d, orderID: %d, totalCost: %s",
			data.AccountId, data.ServiceId, data.OrderId, data.TotalCost.String())

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "given reservation is already cancelled",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	reservation, err = handler.ReservationRepo.GetReservation(tx, data, models.Accepted)
	if err == nil {
		handler.Logger.Errorf("accepted reservation with params: accountID: %d, serviceID: %d, orderID: %d, totalCost: %s",
			data.AccountId, data.ServiceId, data.OrderId, data.TotalCost.String())

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "given reservation was accepted",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	reservation.State = models.Cancelled
	reservation.RecordTime = time.Now()
	err = handler.ReservationRepo.CreateReservation(tx, reservation)
	if err != nil {
		// Shouldn't be there
		return
	}

	handler.TxHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
