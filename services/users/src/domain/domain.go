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
	Msg     string
	Details map[string]any
}

func NewDomainError(ctx context.Context, t DomainErrorType, msg string) *DomainError {
	details, _ := ctx.Value(DomainErrorDetailsKey).(map[string]any)
	return &DomainError{t, msg, details}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf(`%s: %s`, e.Type, e.Msg)
}

type UUID googleUUID.UUID

func NewUUID() UUID {
	return UUID(googleUUID.New())
}

func ParseUUID(ctx context.Context, str string) (UUID, error) {
	uuid, err := googleUUID.Parse(str)
	if err != nil {
		return UUID{}, NewDomainError(ctx, DomainErrorTypeUUIDInvalid, err.Error())
	}
	return UUID(uuid), nil
}

func MustParseUUID(str string) UUID {
	uuid, err := ParseUUID(context.Background(), str)
	if err != nil {
		panic(err)
	}
	return uuid
}

func (uuid UUID) String() string {
	return googleUUID.UUID(uuid).String()
}

func MustRandomString(length int) string {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
