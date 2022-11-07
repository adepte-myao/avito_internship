package handlers

import (
	"encoding/json"
	"net/http"
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

	report, err := handler.Reservation.GetAccountantReport()
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(report)
}
