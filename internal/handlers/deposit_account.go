package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type DepositAccountHandler struct {
	logger      *logrus.Logger
	accountRepo *storage.AccountRepository
	txHelper    *storage.TransactionHelper
}

func NewDepositAccountHandler(logger *logrus.Logger, store *storage.Storage) *DepositAccountHandler {
	return &DepositAccountHandler{
		logger:      logger,
		accountRepo: storage.NewAccountRepository(),
		txHelper:    storage.NewTransactionHelper(store),
	}
}

func (handler *DepositAccountHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Deposit account request received")

	var data dtos.DepositAccountDto
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

	_, err = handler.accountRepo.GetAccount(tx, data.AccountId)
	if err != nil {
		// TODO: can't be other errors except no account?
		err := handler.accountRepo.CreateAccount(tx, data.AccountId)
		if err != nil {
			// TODO
			handler.logger.Error("creation account: ", err.Error())
			return
		}
	}

	err = handler.accountRepo.IncreaseBalance(tx, data.AccountId, data.Value)
	if err != nil {
		// TODO
		handler.logger.Error("increasing balance: : ", err.Error())
		return
	}

	handler.txHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
