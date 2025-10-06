package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"users/src/config"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service interface {
}

type DatabaseRepository interface {
	InsertUser(context.Context, *User) error
	GetUserByID(context.Context, UUID) (*User, error)
}

type CacheRepository interface {
	InsertUser(context.Context, *User) error
	GetUser(context.Context, UUID) (*User, error)

	InsertSession(context.Context, *Session) error
	GetSession(context.Context, UUID) (*Session, error)
}

type service struct {
	config *config.Config

	databaseRepo DatabaseRepository
	cacheRepo    CacheRepository
}

func NewService(
	config *config.Config,
	database DatabaseRepository,
	cache CacheRepository,
) Service {
	return &service{config, database, cache}
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")

	ErrUsernameTaken = errors.New("username taken")
	ErrEmailTaken    = errors.New("email taken")
)

type DomainErrorType string

const (
	DomainErrorTypeUserDoesNotExist    DomainErrorType = "user_does_not_exist"
	DomainErrorTypeSessionDoesNotExist DomainErrorType = "session_does_not_exist"
)

type DomainError struct {
	Type    DomainErrorType
	Message string
	Details map[string]any
}

func newDomainError(ctx context.Context, t DomainErrorType, msg string) *DomainError {
	return &DomainError{t, msg, nil}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf(`%s: %s`, e.Type, e.Message)
}

type UUID struct {
	uuid.UUID
}

func NewUUID() UUID {
	return UUID{uuid.New()}
}

func MustRandomString(length int) string {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
