package api

import (
	"encoding/json"
	"net"
	"net/http"
	"rooms/domain"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (a *api) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func (a *api) handleCreateRoom(w http.ResponseWriter, r *http.Request) {

	// Create a Room for the Session-ID's User.
	err := a.service.CreateRoom(r.Context(), r.FormValue("key"), r.FormValue("name"))
	if err != nil {
		a.httpErr(w, errors.Wrap(err, "cannot create room"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleDeleteRoom(w http.ResponseWriter, r *http.Request) {

	// Delete the Room of the Session-ID's User.
	if err := a.service.DeleteRoom(r.Context(), r.FormValue("key")); err != nil {
		a.httpErr(w, errors.Wrap(err, "cannot delete room"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetRoom(w http.ResponseWriter, r *http.Request) {

	// Get the Room-ID's Room.
	room, err := a.service.GetRoom(r.Context(), mux.Vars(r)["room_id"])
	if err != nil {
		a.httpErr(w, errors.Wrap(err, "cannot get room"))
		return
	}

	// Respond with the Room.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(room); err != nil {
		a.httpErr(w, errors.Wrap(err, "cannot write response"))
	}
}

func (a *api) handleJoinRoom(w http.ResponseWriter, r *http.Request) {

	// Get the room and session IDs from the request's URL.
	urlVals := r.URL.Query()
	sessionID := urlVals.Get("session_id")
	roomID := urlVals.Get("room_id")

	// Have the Session-ID's User join the Room-ID's Room.
	userID, err := a.service.JoinRoom(r.Context(), sessionID, roomID)
	if err != nil {
		a.httpErr(w, errors.Wrap(err, "cannot join room"))
		return
	}

	// Upgrade to a WebSocket connection.
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.httpErr(w, errors.Wrap(err, "cannot upgrade to WebSocket connection"))
		return
	}

	// Create a channel for messages to be writen to the connection.
	sendChan := make(chan domain.ChanMsg[*domain.RoomEvent])

	// Listen to RoomEvents in another go-routine.
	go a.service.ListenForRoomEvents(r.Context(), sendChan, roomID)

	go func() {
		defer conn.Close()

		for {
			select {

			case <-r.Context().Done():

				// Have the Session-ID's User leave the Room-ID's Room.
				if err := a.service.LeaveRoom(r.Context(), userID, roomID); err != nil {
					a.svrErr(errors.Wrap(err, "cannot leave room"))
				}

				return

			// Write ChanMsgs from the send-channel to the connection.
			case msg := <-sendChan:

				// Check if their was an error, and if it was caused by a ProblemDetail.
				pd, ok := errors.Cause(msg.Err).(domain.ProblemDetail)
				if !ok && msg.Err != nil {
					a.wsErr(conn, msg.Err)
					return
				}
				msg.Err = &pd

				// Write the ChanMsg to the connection.
				if err := conn.WriteJSON(msg); err != nil {
					a.wsErr(conn, errors.Wrap(err, "cannot write to WebSocket connection"))
					return
				}
			}
		}
	}()

	for {

		// Listen for messages from the connection as RoomEvents.
		var roomEvent *domain.RoomEvent
		if err := conn.ReadJSON(&roomEvent); err != nil {

			// Check if the error isn't a close error.
			if websocket.IsUnexpectedCloseError(err) || err == net.ErrClosed {
				a.wsErr(conn, errors.Wrap(err, "cannot receive websocket message"))
			}

			conn.Close()
			return
		}

		// Send the RoomEvent for the Room-ID's Room.
		if err := a.service.SendRoomEvent(r.Context(), userID, roomID, roomEvent); err != nil {
			sendChan <- domain.ChanMsg[*domain.RoomEvent]{
				Err: errors.Wrap(err, "cannot send event"),
			}
		}
	}
}
