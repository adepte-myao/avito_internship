package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

// TODO: trunc all digits after dot except two first
type WithdrawAccountHandler struct {
	Logger     *logrus.Logger
	Repository storage.SQLRepository
}

func NewWithdrawAccountHandler(Logger *logrus.Logger, repo storage.SQLRepository) *WithdrawAccountHandler {
	return &WithdrawAccountHandler{
		Logger:     Logger,
		Repository: repo,
	}
}

func (handler *WithdrawAccountHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Withdraw account request received")

	var data dtos.WithdrawAccountDto
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
		// TODO
		handler.Logger.Error("account does not exist")

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "account does not exist",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	if account.Balance.LessThan(data.Value) {
		handler.Logger.Error("Not enough money")

		rw.WriteHeader(http.StatusBadRequest)
		outErr := errors.ResponseError{
			Reason: "not enough money",
		}
		json.NewEncoder(rw).Encode(outErr)
		return
	}

	err = handler.Repository.Account.DecreaseBalance(tx, data.AccountId, data.Value)
	if err != nil {
		// Should not be here
		handler.Logger.Error("decreasing balance: : ", err.Error())
		return
	}

	handler.Repository.SQLTransactionHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
