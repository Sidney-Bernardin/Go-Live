package service

import (
	"context"
	"users/src/domain"

	"github.com/pkg/errors"
)

func (svc *Service) Signup(ctx context.Context, signupForm *domain.SignupForm) (*domain.Session, error) {

	user, err := signupForm.NewUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating user from signup-form")
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

	session := domain.NewSession(user.ID)

	if err = svc.cache.InsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "failed inserting session into cache")
	}

	return session, nil
}
