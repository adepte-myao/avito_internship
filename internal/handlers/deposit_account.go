package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type DepositAccountHandler struct {
	Logger     *logrus.Logger
	Repository storage.SQLRepository
}

func NewDepositAccountHandler(Logger *logrus.Logger, repo storage.SQLRepository) *DepositAccountHandler {
	return &DepositAccountHandler{
		Logger:     Logger,
		Repository: repo,
	}
}

// TODO: trunc all digits after dot except two first
func (handler *DepositAccountHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Deposit account request received")

	var data dtos.DepositAccountDto
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

	_, err = handler.Repository.Account.GetAccount(tx, data.AccountId)
	if err != nil {
		// TODO: can't be other errors except no account?
		err := handler.Repository.Account.CreateAccount(tx, data.AccountId)
		if err != nil {
			// Should not be here
			handler.Logger.Error("creation account: ", err.Error())
			return
		}
	}

	err = handler.Repository.Account.IncreaseBalance(tx, data.AccountId, data.Value)
	if err != nil {
		// Should no be here
		handler.Logger.Error("increasing balance: : ", err.Error())
		return
	}

	handler.Repository.SQLTransactionHelper.CommitTransaction(tx)

	rw.WriteHeader(http.StatusNoContent)
}
