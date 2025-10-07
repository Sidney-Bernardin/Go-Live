package domain

type Session struct {
	ID        UUID
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
	return CSRFToken(MustRandomString(32))
}
