package domain

import (
	"fmt"
)

const (
	PDTypeInvalidInput = "invalid_input"
	PDTypeInvalidID    = "invalid_id"
	PDTypeUnauthorized = "unauthorized"

	PDTypeRoomDoesntExist   = "room_doesnt_exist"
	PDTypeRoomAlreadyExists = "room_already_exists"

	PDTypeViewerDoesntExist   = "viewer_doesnt_exist"
	PDTypeViewerAlreadyExists = "viewer_already_exists"

	PDTypeRoomVideoDoesntExist = "room_video_doesnt_exist"
)

type ProblemDetail struct {
	Type   string `json:"type,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (pd ProblemDetail) Error() string {
	return fmt.Sprintf("%s: '%s'", pd.Type, pd.Detail)
}
