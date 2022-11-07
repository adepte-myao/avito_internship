package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary Withdraw
// @Tags balance
// @Description withdraw account if account exists and has enough money
// @Accept json
// @Produce json
// @Param input body dtos.WithdrawAccountDto true "account id and value to withdraw"
// @Success 204 {integer} integer
// @Failure 400 {object} errors.ResponseError
// @Router /balance/withdraw [post]
func (handler *Handler) withdraw(rw http.ResponseWriter, r *http.Request) {
	// TODO: trunc all digits after dot except two first
	handler.Logger.Info("Withdraw account request received")

	var data dtos.WithdrawAccountDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	err := handler.Account.Withdraw(data.AccountId, data.Value)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
