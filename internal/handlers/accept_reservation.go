package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary AcceptReservation
// @Tags reservation
// @Description accept reserved reservation
// @Accept json
// @Produce json
// @Param input body dtos.ReservationDto true "reservation info"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /reservation/accept [post]
func (handler *Handler) acceptReservation(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Accept reservation request received")

	var data dtos.ReservationDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	err := handler.Reservation.AcceptReservation(data)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
