package domain

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID UUID

	Username     string
	Email        string
	PasswordHash []byte
	PasswordSalt string
}

type SignupForm struct {
	Username string
	Email    string
	Password string
}

func (svc *service) Signup(ctx context.Context, form *SignupForm) (*Session, error) {

	passwordSalt := MustRandomString(32)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password+passwordSalt), 12)
	if err != nil {
		return nil, errors.Wrap(err, "cannot hash passwaord")
	}

	user := &User{
		ID:           NewUUID(),
		Username:     form.Username,
		Email:        form.Email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}

	if err = svc.databaseRepo.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user")
	}

	session := &Session{
		ID:        NewUUID(),
		UserID:    user.ID,
		CSRFToken: MustRandomString(32),
	}

	if err = svc.cacheRepo.InsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "cannot insert session")
	}

	return session, nil
}

func (svc *service) GetUser(ctx context.Context, userID UUID) (*User, error) {

	user, err := svc.cacheRepo.GetUser(ctx, userID)
	if !errors.Is(err, ErrUserNotFound) {
		return user, errors.Wrap(err, "cannot get user from cache")
	}

	user, err = svc.databaseRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, newDomainError(ctx, DomainErrorTypeUserDoesNotExist, "User doesn't exist.")
		}

		return nil, errors.Wrap(err, "cannot get user from database")
	}

	if err := svc.databaseRepo.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user into database")
	}

	return user, nil
}
