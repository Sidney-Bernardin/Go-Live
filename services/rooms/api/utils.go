package api

import (
	"encoding/json"
	"net/http"
	"rooms/domain"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (a *api) svrErr(err error) domain.ProblemDetail {
	a.logger.Error().Stack().Err(err).Msg("Server Error")
	return domain.ProblemDetail{
		Problem: domain.ProblemServerError,
		Detail:  "Server Error",
	}
}

// httpErr writes the cause of e to the connection. If e wasn't caused by a
// ProblemDetail, it's is logged and a new server error ProblemDetail is
// written instead.
func (a *api) httpErr(w http.ResponseWriter, e error) {

	// If the error wasn't caused by a ProblemDetail, treat it as a server error.
	pd, ok := errors.Cause(e).(domain.ProblemDetail)
	if !ok {
		pd = a.svrErr(e)
	}

	// Respond with the ProblemDetail.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(pd.HTTPStatusCode())
	if err := json.NewEncoder(w).Encode(pd); err != nil {
		a.svrErr(errors.Wrap(err, "cannot write response to HTTP connection"))
	}
}

// wsErr writes the cause of e as a close message to the connection. If e wasn't
// caused by a ProblemDetail, it's is logged and a new server error ProblemDetail
// is written instead.
func (a *api) wsErr(conn *websocket.Conn, e error) {
	defer conn.Close()

	// If the error wasn't caused by a ProblemDetail, treat it as a server error.
	pd, ok := errors.Cause(e).(domain.ProblemDetail)
	if !ok {
		pd = a.svrErr(e)
	}

	// Write a close message with the ProblemDetail.
	closeMsg := websocket.FormatCloseMessage(pd.WSCloseCode(), pd.Error())
	err := conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(a.wsCloseTimeout))
	if err != nil {
		a.svrErr(errors.Wrap(err, "cannot write close message to WebSocket connection"))
	}
}
