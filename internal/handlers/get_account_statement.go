package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/errors"
)

// @Summary GetStatement
// @Tags balance
// @Description collects information about all transfers related to given account
// @Accept json
// @Produce json
// @Param input body dtos.GetAccountStatementDto true "account id"
// @Success 200 {object} []models.StatementElem
// @Failure 400 {object} errors.ResponseError
// @Router /balance/statement [get]
func (handler *Handler) getStatement(rw http.ResponseWriter, r *http.Request) {
	handler.Logger.Info("Get account statement request received")

	var data dtos.GetAccountStatementDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handler.Logger.Error("cannot decode request body: ", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(errors.NewErrorInvalidRequestBody(""), rw)
		return
	}

	statements, err := handler.Account.GetStatement(data)
	if err != nil {
		handler.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorToResponse(err, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(statements)
}
