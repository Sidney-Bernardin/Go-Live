package service

import (
	"context"
	"users/src/domain"

	"github.com/pkg/errors"
)

func (svc *Service) GetUserByID(ctx context.Context, userID domain.UUID) (*domain.User, error) {

	user, err := svc.cache.GetUser(ctx, userID)
	if err == nil {
		return user, nil
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, errors.Wrap(err, "failed getting user from cache")
	}

	user, err = svc.database.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &domain.DomainError{
				Type:    domain.DomainErrorTypeUserDoesNotExist,
				Message: "User doesn't exist.",
				Details: map[string]any{
					"user_id": userID.String(),
				},
			}
		}

		return nil, errors.Wrap(err, "failed getting user from database")
	}

	if err := svc.cache.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed inserting user into cache")
	}

	return user, nil
}

func (svc *Service) SearchUsers(ctx context.Context, username domain.Username) ([]*domain.User, error) {
	return nil, nil
}
