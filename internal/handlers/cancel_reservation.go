package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/sirupsen/logrus"
)

type CancelReservationHandler struct {
	Logger  *logrus.Logger
	Service *services.Service
}

func NewCancelReservationHandler(Logger *logrus.Logger, serv *services.Service) *CancelReservationHandler {
	return &CancelReservationHandler{
		Logger:  Logger,
		Service: serv,
	}
}

func (handler *CancelReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Make reservation request received")

	var data dtos.ReservationDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	err := handler.Service.Reservation.CancelReservation(data)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
