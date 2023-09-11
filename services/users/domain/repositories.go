package domain

import (
	"bytes"
	"context"
)

type DatabaseRepository interface {
	CreateAccount(ctx context.Context, profilePicture []byte, user *User) (*LoginResponse, error)
	DeleteAccount(ctx context.Context, userID string) error

	GetUserBySession(ctx context.Context, sessionID string, fields ...string) (*User, error)

	GetUser(ctx context.Context, userID string, fields ...string) (*User, error)
	GetUserByUsername(ctx context.Context, username string, fields ...string) (*User, error)
	SearchUsers(ctx context.Context, username string, offset, limit int64, fields ...string) ([]*User, error)
	UserExists(ctx context.Context, kvPairs ...string) (bool, string, error)

	InsertSession(ctx context.Context, session *Session) (sessionID string, err error)
	DeleteSession(ctx context.Context, sessionID string) error

	GetProfilePicture(ctx context.Context, userID string) (*bytes.Buffer, error)
}
