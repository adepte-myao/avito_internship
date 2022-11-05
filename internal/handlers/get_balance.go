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
	Logger     *logrus.Logger
	Repository storage.SQLRepository
}

func NewGetBalanceHandler(Logger *logrus.Logger, repo storage.SQLRepository) *GetBalanceHandler {
	return &GetBalanceHandler{
		Logger:     Logger,
		Repository: repo,
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

	tx, err := handler.Repository.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		handler.Logger.Error(err.Error())

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer handler.Repository.SQLTransactionHelper.RollbackTransaction(tx)

	account, err := handler.Repository.Account.GetAccount(tx, data.AccountId)
	if err != nil {
		handler.Logger.Error("account with id", data.AccountId, "does not exist")

		rw.WriteHeader(http.StatusBadRequest)
		outError := errors.ResponseError{
			Reason: "account with given ID does not exist",
		}
		json.NewEncoder(rw).Encode(outError)
		return
	}

	handler.Repository.SQLTransactionHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(account)
}
