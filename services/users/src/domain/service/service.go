package service

import (
	"context"
	"users/src"
	"users/src/domain"

	"github.com/pkg/errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")

	ErrUsernameTaken = errors.New("username taken")
	ErrEmailTaken    = errors.New("email taken")
)

type DatabaseRepository interface {
	InsertUser(context.Context, *domain.User) error
	GetUserByID(context.Context, domain.UUID) (*domain.User, error)
}

type CacheRepository interface {
	InsertUser(context.Context, *domain.User) error
	GetUser(context.Context, domain.UUID) (*domain.User, error)

	InsertSession(context.Context, *domain.Session) error
	GetSession(context.Context, domain.UUID) (*domain.Session, error)
}

type Service struct {
	config *src.Config

	databaseRepo DatabaseRepository
	cacheRepo    CacheRepository
}

func NewService(
	config *src.Config,
	databaseRepo DatabaseRepository,
	cacheRepo CacheRepository,
) *Service {
	return &Service{
		config,
		databaseRepo,
		cacheRepo,
	}
}
