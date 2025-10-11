package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	googleUUID "github.com/google/uuid"
)

type contextKey string

type DomainErrorType string

const (
	DomainErrorDetailsKey contextKey = ""

	DomainErrorTypeUserDoesNotExist    DomainErrorType = "user_does_not_exist"
	DomainErrorTypeSessionDoesNotExist DomainErrorType = "session_does_not_exist"

	DomainErrorTypeUUIDInvalid     DomainErrorType = "uuid_invalid"
	DomainErrorTypeUsernameInvalid DomainErrorType = "username_invalid"
	DomainErrorTypePasswordInvalid DomainErrorType = "password_invalid"
)

type DomainError struct {
	Type    DomainErrorType
	Message string
	Details map[string]any
}

func NewDomainError(ctx context.Context, t DomainErrorType, msg string) *DomainError {
	details, _ := ctx.Value(DomainErrorDetailsKey).(map[string]any)
	return &DomainError{t, msg, details}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf(`%s: %s`, e.Type, e.Message)
}

type UUID struct {
	googleUUID.UUID
}

func NewUUID() UUID {
	return UUID{googleUUID.New()}
}

func NewUUIDFromString(ctx context.Context, str string) (UUID, error) {
	uuid, err := googleUUID.Parse(str)
	if err != nil {
		return UUID{}, NewDomainError(ctx, DomainErrorTypeUUIDInvalid, err.Error())
	}
	return UUID{uuid}, nil
}

func MustRandomString(length int) string {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
