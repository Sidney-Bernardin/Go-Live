package service

import (
	"context"
	"users/src/domain"

	"github.com/pkg/errors"
)

func (svc *Service) GetUserByID(ctx context.Context, userID domain.UUID) (*domain.User, error) {

	user, err := svc.cacheRepo.GetUser(ctx, userID)
	if !errors.Is(err, ErrUserNotFound) {
		return user, errors.Wrap(err, "failed getting user from cache")
	}

	user, err = svc.databaseRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &domain.DomainError{
				Type:    domain.DomainErrorTypeUserDoesNotExist,
				Message: "User doesn't exist.",
			}
		}

		return nil, errors.Wrap(err, "failed getting user from database")
	}

	if err := svc.cacheRepo.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed inserting user into cache")
	}

	return user, nil
}

func (svc *Service) SearchUsers(ctx context.Context, username domain.Username) ([]*domain.User, error) {
	return nil, nil
}
