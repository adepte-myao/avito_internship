package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary ProcessInternalTransfer
// @Tags balance
// @Description transfers money between two given accounts
// @Accept json
// @Produce json
// @Param input body dtos.MakeInternalTransferDto true "transfer info"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/transfer [post]
func (handler *Handler) internalTransfer(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Internal transfer request received")

	var data dtos.MakeInternalTransferDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	err := handler.Account.InternalTransfer(data.SenderId, data.ReceiverId, data.Value)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
