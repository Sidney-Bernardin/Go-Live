package http

import (
	"context"
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

// addDataToRequest returns r with the contents of newData added to r's context value.
func (a *api) addDataToRequest(r *http.Request, newData map[string]any) *http.Request {

	// Get request's current data from it's context value.
	currentData, ok := r.Context().Value("data_from_middleware").(map[string]any)
	if !ok {
		currentData = map[string]any{}
	}

	// Add newData to currentData.
	for k, v := range newData {
		currentData[k] = v
	}

	// Update the requests context value with currentData.
	ctx := context.WithValue(r.Context(), "data_from_middleware", currentData)
	return r.WithContext(ctx)
}
