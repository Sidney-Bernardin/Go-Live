package domain

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	ProblemServerError  = "server_error"
	ProblemUnauthorized = "unauthorized"

	ProblemInvalidInput = "invalid_input"
	ProblemInvalidID    = "invalid_id"

	ProblemRoomAlreadyExists   = "room_already_exists"
	ProblemViewerAlreadyExists = "viewer_already_exists"

	ProblemRoomDoesntExist   = "room_doesnt_exist"
	ProblemViewerDoesntExist = "viewer_doesnt_exist"
)

type ProblemDetail struct {
	Problem string `json:"problem,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (pd ProblemDetail) Error() string {
	return fmt.Sprintf("%s: %s", pd.Problem, pd.Detail)
}

func (pd ProblemDetail) HTTPStatusCode() int {
	switch pd.Problem {
	case ProblemServerError:
		return http.StatusInternalServerError
	case ProblemUnauthorized:
		return http.StatusUnauthorized
	case ProblemInvalidInput, ProblemInvalidID:
		return http.StatusUnprocessableEntity
	case ProblemRoomDoesntExist, ProblemViewerDoesntExist:
		return http.StatusNotFound
	case ProblemRoomAlreadyExists, ProblemViewerAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

func (pd ProblemDetail) WSCloseCode() int {
	switch pd.Problem {
	case ProblemServerError:
		return websocket.CloseInternalServerErr
	default:
		return websocket.CloseAbnormalClosure
	}
}
