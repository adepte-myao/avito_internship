package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/sirupsen/logrus"
)

// TODO: trunc all digits after dot except two first
type WithdrawAccountHandler struct {
	Logger  *logrus.Logger
	Service *services.Service
}

func NewWithdrawAccountHandler(Logger *logrus.Logger, serv *services.Service) *WithdrawAccountHandler {
	return &WithdrawAccountHandler{
		Logger:  Logger,
		Service: serv,
	}
}

func (handler *WithdrawAccountHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Withdraw account request received")

	var data dtos.WithdrawAccountDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	err := handler.Service.Account.Withdraw(data)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
