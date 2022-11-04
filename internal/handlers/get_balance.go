package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type GetBalanceHandler struct {
	logger      *logrus.Logger
	accountRepo *storage.AccountRepository
	txHelper    *storage.TransactionHelper
}

func NewGetBalanceHandler(logger *logrus.Logger, store *storage.Storage) *GetBalanceHandler {
	return &GetBalanceHandler{
		logger:      logger,
		accountRepo: storage.NewAccountRepository(),
		txHelper:    storage.NewTransactionHelper(store),
	}
}

func (handler *GetBalanceHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Get balance request received")

	var data dtos.GetBalanceDto
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
		handler.logger.Error("Account does not exist")
		return
	}

	handler.txHelper.CommitTransaction(tx)

	err = json.NewEncoder(rw).Encode(account)
	if err != nil {
		// TODO
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
