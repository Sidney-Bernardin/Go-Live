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
