package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary GetAccountantReport
// @Tags accountant
// @Description collects information about accepted reservations
// @Produce json
// @Success 200 {object} []models.AccountantReportElem
// @Failure 400 {object} errors.ResponseError
// @Router /accountant-report [get]
func (handler *Handler) getAccountantReport(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Get accountant report request received")

	var data dtos.MakeAccountantReportDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	report, err := handler.Reservation.GetAccountantReport(data.Month, data.Year)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(report)
}
