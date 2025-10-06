package domain

type Session struct {
	ID        UUID
	UserID    UUID
	CSRFToken string
}
