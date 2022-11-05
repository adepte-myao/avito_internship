package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/sirupsen/logrus"
)

type GetBalanceHandler struct {
	Logger  *logrus.Logger
	Service *services.Service
}

func NewGetBalanceHandler(Logger *logrus.Logger, serv *services.Service) *GetBalanceHandler {
	return &GetBalanceHandler{
		Logger:  Logger,
		Service: serv,
	}
}

func (handler *GetBalanceHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Get balance request received")

	var data dtos.GetBalanceDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	account, err := handler.Service.Account.GetBalance(data)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(account)
}
