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
	Logger     *logrus.Logger
	Repository storage.SQLRepository
}

func NewMakeReservationHandler(Logger *logrus.Logger, repo storage.SQLRepository) *MakeReservationHandler {
	return &MakeReservationHandler{
		Logger:     Logger,
		Repository: repo,
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

	tx, err := handler.Repository.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		// Should not be here
		return
	}
	defer handler.Repository.SQLTransactionHelper.RollbackTransaction(tx)

	account, err := handler.Repository.Account.GetAccount(tx, data.AccountId)
	if err != nil {
		handler.Logger.Errorf("account with ID %d does not exist", data.AccountId)

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "account does not exist",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	if account.Balance.LessThan(data.TotalCost) {
		handler.Logger.Errorf("not enough money: account: id: %d, balance: %s; required: %s",
			account.ID, account.Balance.String(), data.TotalCost.String())

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "not enough money",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	err = handler.Repository.Account.DecreaseBalance(tx, data.AccountId, data.TotalCost)
	if err != nil {
		// Should not be here
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
	err = handler.Repository.Reservation.CreateReservation(tx, reservation)
	if err != nil {
		// Should not be here
		return
	}

	handler.Repository.SQLTransactionHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
