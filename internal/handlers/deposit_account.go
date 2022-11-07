package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary Deposit
// @Tags balance
// @Description deposit account if account exists, otherwise creates an account and deposi it
// @Accept json
// @Produce json
// @Param input body dtos.DepositAccountDto true "account id and value to deposit"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/deposit [post]
func (handler *Handler) deposit(rw http.ResponseWriter, r *http.Request) {
	// TODO: trunc all digits after dot except two first
	handler.Logger.Info("Deposit account request received")

	var data dtos.DepositAccountDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	err := handler.Account.Deposit(data.AccountId, data.Value)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
