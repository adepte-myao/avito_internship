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
}

func NewGetBalanceHandler(logger *logrus.Logger, store *storage.Storage) *GetBalanceHandler {
	return &GetBalanceHandler{
		logger:      logger,
		accountRepo: storage.NewAccountRepository(store),
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

	account, err := handler.accountRepo.GetAccount(data.AccountId)
	if err != nil {
		handler.logger.Error("Account does not exist")
		return
	}

	err = json.NewEncoder(rw).Encode(account)
	if err != nil {
		// TODO
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
