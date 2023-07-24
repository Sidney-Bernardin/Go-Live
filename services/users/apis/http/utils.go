package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// err writes e to the given http.ResponseWriter with with the given status code.
func (a *api) err(w http.ResponseWriter, statusCode int, e error) {

	// If the status code 500 or more, treat e as a server error.
	if statusCode >= 500 {
		a.logger.Error().Stack().Err(e).Msg("Server Error")
		e = errors.New("Server Error")
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		err = errors.Wrap(err, "cannot write response to HTTP connection")
		a.logger.Error().Stack().Err(err).Msg("Server Error")
	}
}
