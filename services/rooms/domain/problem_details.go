package domain

import "fmt"

const (
	PDTypeUnauthorized = "unauthorized"

	PDTypeInvalidInput = "invalid_input"
	PDTypeBadStreamID  = "bad_stream_id"
	PDTypeInvalidID    = "invalid_id"

	PDTypeRoomDoesntExist   = "room_doesnt_exist"
	PDTypeRoomAlreadyExists = "room_already_exists"

	PDTypeViewerDoesntExist   = "viewer_doesnt_exist"
	PDTypeViewerAlreadyExists = "viewer_already_exists"
)

type ProblemDetail struct {
	Type   string `json:"type,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (pd ProblemDetail) Error() string {
	return fmt.Sprintf("%s: '%s'", pd.Type, pd.Detail)
}
