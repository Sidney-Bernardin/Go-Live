package domain

import (
	"fmt"
)

type DomainErrorType string

const (
	DomainErrorTypeUserDoesNotExist    DomainErrorType = "user_does_not_exist"
	DomainErrorTypeSessionDoesNotExist DomainErrorType = "session_does_not_exist"

	DomainErrorTypeUUIDInvalid           DomainErrorType = "uuid_invalid"
	DomainErrorTypeUsernameInvalid       DomainErrorType = "username_invalid"
	DomainErrorTypePasswordInvalid       DomainErrorType = "password_invalid"
	DomainErrorTypeEmailInvalid          DomainErrorType = "email_invalid"
	DomainErrorTypeProfilePictureInvalid DomainErrorType = "profile_picture_invalid"

	DomainErrorTypeAuthenticationFailed DomainErrorType = "authentication_failed"
)

type DomainError struct {
	Type    DomainErrorType
	Message string
	Details map[string]any
}

func (e *DomainError) Error() string {
	return fmt.Sprintf(`%s: %s`, e.Type, e.Message)
}
