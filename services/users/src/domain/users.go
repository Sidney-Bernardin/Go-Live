package domain

import (
	"time"
	"users/src"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Username     Username
	Email        string
	PasswordHash PasswordHash
	PasswordSalt PasswordSalt
}

type Username string

func NewUsername(u string) (Username, error) {
	if len(u) < 3 || 32 < len(u) {
		return "", &DomainError{DomainErrorTypeUsernameInvalid, "Username must be between 3 and 32 characters.", map[string]any{
			"username": u,
		}}
	}
	return Username(u), nil
}

type Password string

func NewPassword(pw string) (Password, error) {
	if len(pw) < 8 || 100 < len(pw) {
		return "", &DomainError{DomainErrorTypePasswordInvalid, "Password must be between 8 and 100 characters.", map[string]any{
			"password": pw,
		}}
	}
	return Password(pw), nil
}

type PasswordSalt string

func NewPasswordSalt() PasswordSalt {
	return PasswordSalt(src.MustRandomString(32))
}

type PasswordHash []byte

func NewPasswordHash(pw Password, pwSalt PasswordSalt) (PasswordHash, error) {
	pwHash := []byte(string(pw) + string(pwSalt))
	pwHash, err := bcrypt.GenerateFromPassword(pwHash, 12)
	if err != nil {
		return nil, errors.Wrap(err, "failed hashing passwaord")
	}
	return PasswordHash(pwHash), nil
}
