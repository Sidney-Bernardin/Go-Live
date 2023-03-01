package domain

import "bytes"

type DatabaseRepository interface {
	InsertUserWithSession(user *User, session *Session) (sessionID string, err error)
	GetUser(userID string, fields ...string) (*User, error)
	SearchUsers(username string, offset, limit int64, fields ...string) ([]*User, error)
	GetUserByUsername(username string, fields ...string) (*User, error)
	CheckForTakenUserFields(fields map[string]any) error

	GetSessionsUser(sessionID string, fields ...string) (*User, error)

	InsertSession(session *Session) (sessionID string, err error)
	DeleteSession(sessionID string) error

	InsertProfilePicture(sessionID string, fileBytes []byte) error
	GetProfilePicture(userID string) (*bytes.Buffer, error)
}
