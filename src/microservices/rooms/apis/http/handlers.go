package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rooms/domain"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (a *api) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func (a *api) handleCreateRoom(w http.ResponseWriter, r *http.Request) {

	roomSettings := &domain.RoomSettings{"No Name"}

	if err := a.service.CreateRoom(r.FormValue("key"), r.FormValue("name"), roomSettings); err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot create room"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleDeleteRoom(w http.ResponseWriter, r *http.Request) {

	if err := a.service.DeleteRoom(r.FormValue("key")); err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot delete room"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleAuthenticateStream(w http.ResponseWriter, r *http.Request) {

	_ = r.Header.Get("X-Original-Uri")

	// Respond.
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetRoom(w http.ResponseWriter, r *http.Request) {

	room, err := a.service.GetRoom(mux.Vars(r)["room_id"])
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

func (a *api) handleJoinRoom(w http.ResponseWriter, r *http.Request) {

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

	sessionID := r.URL.Query().Get("session_id")
	roomID := r.URL.Query().Get("room_id")

	user, roomChan, err := a.service.JoinRoom(sessionID, roomID)
	if err != nil {

		// If err was caused by a problem-detail, close with AbnormalClosure.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.wsCloseErr(conn, websocket.CloseAbnormalClosure, pd)
			return
		}

		// Close with InternalServerErr.
		a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot join room"))
		return
	}

	// In another go-routine, listen to the room's channel and write it's
	// messages to the client.
	go func() {

		defer func() {
			conn.Close()
			if err := a.service.LeaveRoom(user.ID, roomID); err != nil {
				a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot leave room"))
			}
		}()

		for {
			select {

			case <-r.Context().Done():
				return

			case msg, ok := <-roomChan:

				// Check if the room's channel was closed.
				if !ok {
					return
				}

				// Write the message to the client.
				if err := conn.WriteJSON(msg); err != nil {

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
		_, msgBytes, err := conn.ReadMessage()
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
		var msg map[string]any
		if err := json.Unmarshal(msgBytes, &msg); err != nil {

			// Close with AbnormalClosure.
			a.wsCloseErr(conn, websocket.CloseAbnormalClosure, domain.ProblemDetail{
				Type:   PDTypeCannotProcessRequestData,
				Detail: fmt.Sprintf("Cannot decode message: '%v'", err),
			})

			return
		}

		if err := a.service.BroadcastMessage(user, roomID, msg); err != nil {

			// If err was caused by a problem-detail, close with AbnormalClosure.
			if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
				a.wsCloseErr(conn, websocket.CloseAbnormalClosure, pd)
				return
			}

			// Close with InternalServerErr.
			a.wsCloseErr(conn, websocket.CloseInternalServerErr, errors.Wrap(err, "cannot broadcast message"))
			return
		}
	}
}
