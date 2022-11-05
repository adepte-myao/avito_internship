package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/errors"
)

func writeErrorToResponse(err error, rw http.ResponseWriter) {
	outErr := errors.ResponseError{
		Reason: err.Error(),
	}
	json.NewEncoder(rw).Encode(outErr)
}
