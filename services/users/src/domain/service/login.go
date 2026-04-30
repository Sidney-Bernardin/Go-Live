package service

import (
	"context"
	"io"
	"users/src"
	"users/src/domain"

	"github.com/pkg/errors"
)

type SignupForm struct {
	Username       string
	Email          string
	Password       string
	ProfilePicture io.Reader
}

func (svc *Service) Signup(ctx context.Context, form *SignupForm) (*domain.Session, error) {

	username, err := domain.NewUsername(form.Username)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating username")
	}

	password, err := domain.NewPassword(form.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating password")
	}

	passwordSalt := domain.NewPasswordSalt()
	passwordHash, err := domain.NewPasswordHash(password, passwordSalt)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating passwaord-hash")
	}

	profilePicture, err := domain.NewProfilePicture(svc.config, form.ProfilePicture)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating profile-picture")
	}

	user := &domain.User{
		ID:           domain.NewUUID(),
		Username:     username,
		Email:        form.Email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}

	if err = svc.database.InsertUser(ctx, user); err != nil {
		if errors.Is(err, ErrUsernameTaken) {
			return nil, &domain.DomainError{
				Type:    domain.DomainErrorTypeUsernameInvalid,
				Message: "Username already in use.",
				Details: map[string]any{
					"username": user.Username,
				},
			}
		}

		if errors.Is(err, ErrEmailTaken) {
			return nil, &domain.DomainError{
				Type:    domain.DomainErrorTypeEmailInvalid,
				Message: "Email already in use.",
				Details: map[string]any{
					"email": user.Email,
				},
			}
		}

		return nil, errors.Wrap(err, "failed inserting user into database")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), svc.config.ProfilePictureInsertTimeout)
		defer cancel()
		if err := svc.blobStore.InsertProfilePicture(ctx, user.ID, profilePicture); err != nil {
			svc.logger.Error("Failed to insert profile-picture.", src.ErrGroup(err))
		}
	}()

	session := domain.NewSession(user.ID)

	if err = svc.cache.InsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "failed inserting session into cache")
	}

	return session, nil
}

type SigninForm struct {
	Username string
	Password string
}

func (svc *Service) Signin(ctx context.Context, form *SigninForm) (*domain.Session, error) {

	errUnauthenticated := &domain.DomainError{
		Type:    domain.DomainErrorTypeAuthenticationFailed,
		Message: "Invalid username or password.",
		Details: map[string]any{
			"username": form.Username,
			"password": form.Password,
		},
	}

	username, err := domain.NewUsername(form.Username)
	if err != nil {
		return nil, errUnauthenticated
	}

	password, err := domain.NewPassword(form.Password)
	if err != nil {
		return nil, errUnauthenticated
	}

	user, err := svc.database.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errUnauthenticated
		}

		return nil, errors.Wrap(err, "failed getting user from database")
	}

	if !user.PasswordHash.Compare(password, user.PasswordSalt) {
		return nil, errUnauthenticated
	}

	session := domain.NewSession(user.ID)

	if err = svc.cache.InsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "failed inserting session into cache")
	}

	return session, nil
}
