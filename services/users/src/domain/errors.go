package domain

import (
	"context"
	"fmt"
)

type DomainErrorType string

const (
	DomainErrorTypeUserDoesNotExist    DomainErrorType = "user_does_not_exist"
	DomainErrorTypeSessionDoesNotExist DomainErrorType = "session_does_not_exist"

	DomainErrorTypeUsernameInvalid DomainErrorType = "username_invalid"
	DomainErrorTypePasswordInvalid DomainErrorType = "password_invalid"
)

type DomainError struct {
	Type    DomainErrorType
	Message string
	Details map[string]any
}

func NewDomainError(ctx context.Context, t DomainErrorType, msg string) *DomainError {
	return &DomainError{t, msg, nil}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf(`%s: %s`, e.Type, e.Message)
}
