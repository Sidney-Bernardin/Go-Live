package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rooms/domain"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (a *api) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func (a *api) handleCallbacks(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	callback := mwData["callback"].(string)
	sessionID := mwData["session_id"].(string)
	name := mwData["name"].(string)

	var err error

	switch callback {
	case "create":
		err = a.service.CreateRoom(sessionID, &domain.RoomSettings{
			Name: name,
		})

	case "delete":
		err = a.service.DeleteRoom(sessionID)
	case "join":
		err = a.service.JoinRoom(sessionID, name)
	case "leave":
		err = a.service.LeaveRoom(sessionID, name)
	default:

		// Respond with StatusBadRequest.
		a.err(w, http.StatusBadRequest, domain.ProblemDetail{
			Type: PDTypeInvalidCallback,
		})

		return
	}

	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrapf(err, "cannot %s room", callback))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetRoom(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	roomID := mwData["room_id"].(string)

	room, err := a.service.GetRoom(roomID)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get video"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(room); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleDiagnostics(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	sessionID := mwData["session_id"].(string)
	roomID := mwData["room_id"].(string)

	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {

		// Respond with a StatusBadRequest.
		a.err(w, http.StatusBadRequest, domain.ProblemDetail{
			Type:   PDTypeCannotUpgradeRequest,
			Detail: fmt.Sprintf("Cannot upgrade HTTP connection to a WebSocket connection: '%v'", err),
		})
	}
	defer conn.Close()

	user, diagnosticsChan, err := a.service.NewDiagnosticsListener(sessionID, roomID)
	if err != nil {

		// If err was caused by a problem-detail, close with AbnormalClosure.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.wsCloseErr(conn, websocket.CloseAbnormalClosure, pd)
			return
		}

		// Close with InternalServerErr.
		a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot get a diagnostics listener"))
		return
	}

	// In another go-routine, listen to the diagnostics-channel and write the
	// diagnostics to the client.
	go func() {

		defer conn.Close()

		for {
			select {

			case <-r.Context().Done():
				return

			case diagnostic, ok := <-diagnosticsChan:

				// Check if the channel was closed.
				if !ok {
					return
				}

				// Write the diagnostic to the client.
				if err := conn.WriteJSON(diagnostic); err != nil {

					// Check for close errors.
					_, isCloseErr := err.(*websocket.CloseError)
					if isCloseErr || strings.HasSuffix(err.Error(), ": use of closed network connection") || err == websocket.ErrCloseSent {
						return
					}

					// Close with InternalServerErr.
					a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot write to WebSocket connection"))
					return
				}
			}
		}
	}()

	for {

		// Listen for messages from the client.
		_, msg, err := conn.ReadMessage()
		if err != nil {

			// Check for close errors.
			_, isCloseErr := err.(*websocket.CloseError)
			if isCloseErr || strings.HasSuffix(err.Error(), ": use of closed network connection") || err == websocket.ErrCloseSent {
				return
			}

			// Close with InternalServerErr.
			a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot listen to WebSocket connection"))
			return
		}

		// Decode the message.
		var diagnostic domain.Diagnostic
		if err := json.Unmarshal(msg, &diagnostic); err != nil {

			// Close with AbnormalClosure.
			a.wsCloseErr(conn, websocket.CloseAbnormalClosure, domain.ProblemDetail{
				Type:   PDTypeCannotProcessRequestData,
				Detail: fmt.Sprintf("Cannot decode message: '%v'", err),
			})

			return
		}

		if err := a.service.BroadcastDiagnostic(user, roomID, &diagnostic); err != nil {

			// If err was caused by a problem-detail, close with AbnormalClosure.
			if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
				a.wsCloseErr(conn, websocket.CloseAbnormalClosure, pd)
				return
			}

			// Close with InternalServerErr.
			a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot broadcast diagnostic"))
			return
		}
	}
}
