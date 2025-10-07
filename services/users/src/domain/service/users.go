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

	if err = svc.databaseRepo.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed inserting user into database")
	}

	session := domain.NewSession(user.ID)

	if err = svc.cacheRepo.InsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "failed inserting session into cache")
	}

	return session, nil
}

func (svc *Service) GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error) {

	user, err := svc.cacheRepo.GetUser(ctx, userID)
	if !errors.Is(err, ErrUserNotFound) {
		return user, errors.Wrap(err, "failed getting user from cache")
	}

	user, err = svc.databaseRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, domain.NewDomainError(ctx, domain.DomainErrorTypeUserDoesNotExist, "User doesn't exist.")
		}

		return nil, errors.Wrap(err, "failed getting user from database")
	}

	if err := svc.cacheRepo.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed inserting user into cache")
	}

	return user, nil
}
