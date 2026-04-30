package service

import (
	"context"
	"log/slog"
	"users/src"
	"users/src/domain"

	"github.com/pkg/errors"
)

type Service struct {
	config *src.Config
	logger *slog.Logger

	database  DatabaseRepository
	cache     CacheRepository
	blobStore BlobStoreRepository
}

type (
	DatabaseRepository interface {
		InsertUser(ctx context.Context, user *domain.User) error
		GetUserByID(ctx context.Context, userID domain.UUID) (*domain.User, error)
		GetUserByUsername(ctx context.Context, username domain.Username) (*domain.User, error)
	}

	CacheRepository interface {
		InsertUser(ctx context.Context, user *domain.User) error
		GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error)

		InsertSession(ctx context.Context, session *domain.Session) error
		GetSession(ctx context.Context, sessionID domain.UUID) (*domain.Session, error)
	}

	BlobStoreRepository interface {
		InsertProfilePicture(ctx context.Context, userID domain.UUID, profilePicture domain.ProfilePicture) error
	}
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")

	ErrUsernameTaken = errors.New("username taken")
	ErrEmailTaken    = errors.New("email taken")
)

func New(
	config *src.Config,
	logger *slog.Logger,
	database DatabaseRepository,
	cache CacheRepository,
	blobStore BlobStoreRepository,
) *Service {
	return &Service{
		config,
		logger,
		database,
		cache,
		blobStore,
	}
}
