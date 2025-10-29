package domain

import (
	"context"

	"github.com/pkg/errors"
)

type SignupForm struct {
	Username string
	Email    string
	Password string
}

func (signupForm *SignupForm) NewUser(ctx context.Context) (*User, error) {

	username, err := NewUsername(signupForm.Username)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating username")
	}

	password, err := NewPassword(signupForm.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating password")
	}

	passwordSalt := NewPasswordSalt()
	passwordHash, err := NewPasswordHash(password, passwordSalt)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating passwaord-hash")
	}

	return &User{
		ID:           NewUUID(),
		Username:     username,
		Email:        signupForm.Email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}, nil
}
