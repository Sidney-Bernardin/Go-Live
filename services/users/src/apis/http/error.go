package http

import (
	"fmt"
	"net/http"
	"users/src/domain"
)

var apiCodes = map[apiErrorType]int{
	apiErrorTypeInternalServerError:   http.StatusInternalServerError,
	apiErrorTypeUnauthorized:          http.StatusUnauthorized,
	apiErrorTypeProfilePictureInvalid: http.StatusBadRequest,
}

var domainCodes = map[domain.DomainErrorType]int{
	domain.DomainErrorTypeUserDoesNotExist:    http.StatusNotFound,
	domain.DomainErrorTypeSessionDoesNotExist: http.StatusNotFound,

	domain.DomainErrorTypeUUIDInvalid:           http.StatusBadRequest,
	domain.DomainErrorTypeUsernameInvalid:       http.StatusBadRequest,
	domain.DomainErrorTypePasswordInvalid:       http.StatusBadRequest,
	domain.DomainErrorTypeEmailInvalid:          http.StatusBadRequest,
	domain.DomainErrorTypeProfilePictureInvalid: http.StatusBadRequest,

	domain.DomainErrorTypeAuthenticationFailed: http.StatusUnauthorized,
}

type apiErrorType string

const (
	apiErrorTypeInternalServerError   apiErrorType = "internal_server_error"
	apiErrorTypeUnauthorized          apiErrorType = "unauthorized"
	apiErrorTypeProfilePictureInvalid apiErrorType = "profile_picture_invalid"
)

type apiError[T apiErrorType | domain.DomainErrorType] struct {
	Type    T              `json:"type"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *apiError[T]) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
