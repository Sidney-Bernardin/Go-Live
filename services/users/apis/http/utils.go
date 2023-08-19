package http

import (
	"encoding/json"
	"net/http"
	"users/domain"

	"github.com/pkg/errors"
)

// err writes the cause of e to the connection. If e wasn't caused by a
// ProblemDetail, it's is logged and a new server error ProblemDetail is
// written instead.
func (a *api) err(w http.ResponseWriter, e error) {

	// If the error wasn't caused by a ProblemDetail, treat it as a server error.
	pd, ok := errors.Cause(e).(domain.ProblemDetail)
	if !ok {
		a.logger.Error().Stack().Err(e).Msg("Server Error")
		pd = domain.ProblemDetail{
			Problem: domain.ProblemServerError,
			Detail:  "Server Error",
		}
	}

	// Respond with the ProblemDetail.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(pd.HTTPStatusCode())
	if err := json.NewEncoder(w).Encode(pd); err != nil {
		err = errors.Wrap(err, "cannot write response to HTTP connection")
		a.logger.Error().Stack().Err(err).Msg("Server Error")
	}
}
