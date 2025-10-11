package domain

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID UUID

	Username     Username
	Email        string
	PasswordHash PasswordHash
	PasswordSalt PasswordSalt
}

type Username string

func NewUsername(ctx context.Context, uname string) (Username, error) {
	if len(uname) < 3 || 32 < len(uname) {
		return "", NewDomainError(ctx, DomainErrorTypeUsernameInvalid, "Username must be between 3 and 32 characters.")
	}
	return Username(uname), nil
}

type Password string

func NewPassword(ctx context.Context, passw string) (Password, error) {
	if len(passw) < 8 || 100 < len(passw) {
		return "", NewDomainError(ctx, DomainErrorTypePasswordInvalid, "Password must be between 8 and 100 characters.")
	}
	return Password(passw), nil
}

type PasswordSalt string

func NewPasswordSalt() PasswordSalt {
	return PasswordSalt(MustRandomString(32))
}

type PasswordHash []byte

func NewPasswordHash(ctx context.Context, password Password, passwordSalt PasswordSalt) (PasswordHash, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(string(password)+string(passwordSalt)), 12)
	if err != nil {
		return nil, errors.Wrap(err, "failed hashing passwaord")
	}
	return PasswordHash(passwordHash), nil
}
