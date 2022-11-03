package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/sirupsen/logrus"
)

type AcceptReservationHandler struct {
	logger          *logrus.Logger
	accountRepo     *storage.AccountRepository
	reservationRepo *storage.ReservationRepository
}

func NewAcceptReservationHandler(logger *logrus.Logger, store *storage.Storage) *AcceptReservationHandler {
	return &AcceptReservationHandler{
		logger:          logger,
		accountRepo:     storage.NewAccountRepository(store),
		reservationRepo: storage.NewReservationRepository(store),
	}
}

func (handler *AcceptReservationHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Make reservation request received")

	var data dtos.ReservationDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.logger.Error("cannot decode request body: ", err.Error())
		// TODO: find out what part of body was not decoded, maybe get more user-friendly output
		return
	}

	reservation, err := handler.reservationRepo.GetReservation(data)
	if err != nil {
		// TODO
		return
	}

	reservation.State = models.Accepted
	reservation.RecordTime = time.Now()
	err = handler.reservationRepo.CreateReservation(reservation)
	if err != nil {
		// TODO
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
