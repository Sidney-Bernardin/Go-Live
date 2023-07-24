package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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

// err writes e to the given WebSocket connection with with the given close code.
func (a *api) wsCloseErr(conn *websocket.Conn, closeCode int, e error) {

	// If the close code is 1011, treat e as a server error.
	if closeCode == websocket.CloseInternalServerErr {
		a.logger.Error().Stack().Err(e).Msg("Server Error")
		e = errors.New("Server Error")
	}

	// Write the error as a control message.
	b := websocket.FormatCloseMessage(closeCode, e.Error())
	if err := conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(a.pongTimeout)); err != nil {

		// Check for close errors.
		_, isCloseErr := err.(*websocket.CloseError)
		if isCloseErr || strings.HasSuffix(err.Error(), ": use of closed network connection") || err == websocket.ErrCloseSent {
			return
		}

		err = errors.Wrap(err, "cannot write close message to WebSocket connection")
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
