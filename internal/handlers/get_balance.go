package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type GetBalanceHandler struct {
	Logger      *logrus.Logger
	AccountRepo storage.AccountRepo
	TxHelper    storage.SQLTransactionHelper
}

func NewGetBalanceHandler(Logger *logrus.Logger, store *storage.Storage) *GetBalanceHandler {
	return &GetBalanceHandler{
		Logger:      Logger,
		AccountRepo: storage.NewAccountRepository(),
		TxHelper:    storage.NewTransactionHelper(store),
	}
}

func (handler *GetBalanceHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Get balance request received")

	var data dtos.GetBalanceDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())

		rw.WriteHeader(http.StatusBadRequest)
		outError := errors.ResponseError{
			Reason: "invalid request body",
		}
		json.NewEncoder(rw).Encode(outError)
		return
	}

	tx, err := handler.TxHelper.BeginTransaction()
	if err != nil {
		handler.Logger.Error(err.Error())

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer handler.TxHelper.RollbackTransaction(tx)

	account, err := handler.AccountRepo.GetAccount(tx, data.AccountId)
	if err != nil {
		handler.Logger.Error("account with id", data.AccountId, "does not exist")

		rw.WriteHeader(http.StatusBadRequest)
		outError := errors.ResponseError{
			Reason: "account with given ID does not exist",
		}
		json.NewEncoder(rw).Encode(outError)
		return
	}

	handler.TxHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(account)
}
