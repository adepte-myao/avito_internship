package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary GetBalance
// @Tags balance
// @Description return balance of the account with given id
// @Accept json
// @Produce json
// @Param input body dtos.GetBalanceDto true "account id"
// @Success 200 {object} models.Account
// @Failure 400 {object} errors.ResponseError
// @Router /balance/get [get]
func (handler *Handler) getBalance(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Get balance request received")

	var data dtos.GetBalanceDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	account, err := handler.Account.GetBalance(data.AccountId)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(account)
}
