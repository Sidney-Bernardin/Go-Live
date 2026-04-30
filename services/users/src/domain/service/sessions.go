package service

import (
	"context"
	"users/src/domain"

	"github.com/pkg/errors"
)

func (svc *Service) Authenticate(ctx context.Context, sessionID domain.UUID) (*domain.Session, error) {

	session, err := svc.cache.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, &domain.DomainError{
				Type: domain.DomainErrorTypeAuthenticationFailed,
				Details: map[string]any{
					"session_id": sessionID.String(),
				},
			}
		}

		return nil, errors.Wrap(err, "failed getting session")
	}

	return session, nil
}
