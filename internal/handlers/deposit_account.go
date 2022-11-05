package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/sirupsen/logrus"
)

type DepositAccountHandler struct {
	Logger  *logrus.Logger
	Service *services.Service
}

func NewDepositAccountHandler(Logger *logrus.Logger, serv *services.Service) *DepositAccountHandler {
	return &DepositAccountHandler{
		Logger:  Logger,
		Service: serv,
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

	err := handler.Service.Account.Deposit(data)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
