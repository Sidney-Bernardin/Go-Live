package domain

import "users/src"

type Session struct {
	ID UUID

	UserID    UUID
	CSRFToken CSRFToken
}

func NewSession(userID UUID) *Session {
	return &Session{
		ID:        NewUUID(),
		UserID:    userID,
		CSRFToken: NewCSRFToken(),
	}
}

type CSRFToken string

func NewCSRFToken() CSRFToken {
	return CSRFToken(src.MustRandomString(32))
}
