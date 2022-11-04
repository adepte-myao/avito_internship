package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type WithdrawAccountHandler struct {
	logger      *logrus.Logger
	accountRepo *storage.AccountRepository
	txHelper    *storage.TransactionHelper
}

func NewWithdrawAccountHandler(logger *logrus.Logger, store *storage.Storage) *WithdrawAccountHandler {
	return &WithdrawAccountHandler{
		logger:      logger,
		accountRepo: storage.NewAccountRepository(),
		txHelper:    storage.NewTransactionHelper(store),
	}
}

func (handler *WithdrawAccountHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Withdraw account request received")

	var data dtos.WithdrawAccountDto
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
		handler.logger.Error("Account does not exist")
		return
	}
	if account.Balance.LessThan(data.Value) {
		handler.logger.Error("Not enough money")
		return
	}

	err = handler.accountRepo.DecreaseBalance(tx, data.AccountId, data.Value)
	if err != nil {
		// TODO
		handler.logger.Error("decreasing balance: : ", err.Error())
		return
	}

	handler.txHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
