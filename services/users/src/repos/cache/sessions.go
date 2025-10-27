package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type repoSession struct {
	ID uuid.UUID `json:"id"`

	UserID    uuid.UUID `json:"user_id"`
	CSRFToken string    `json:"csrf_token"`
}

func (s *repoSession) domainify() *domain.Session {
	if s == nil {
		return nil
	}

	return &domain.Session{
		ID:        domain.UUID(s.ID),
		UserID:    domain.UUID(s.UserID),
		CSRFToken: domain.CSRFToken(s.CSRFToken),
	}
}

func (repo *repository) InsertSession(ctx context.Context, session *domain.Session) error {
	key := fmt.Sprintf("session:%s", session.ID)
	p := repo.client.TxPipeline()

	// Set the session.
	jsonSetCmd := p.JSONSet(ctx, key, ".", &repoSession{
		ID:        uuid.UUID(session.ID),
		UserID:    uuid.UUID(session.UserID),
		CSRFToken: string(session.CSRFToken),
	})

	// Set the session's expiry.
	expireCmd := p.Expire(ctx, key, repo.config.SessionDuration)

	if _, err := p.Exec(ctx); err != nil {
		return errors.Wrap(err, "command executions failed")
	}

	if err := jsonSetCmd.Err(); err != nil {
		return errors.Wrap(err, "cannot set json")
	}

	if err := expireCmd.Err(); err != nil {
		return errors.Wrap(err, "cannot set expiry")
	}

	return nil
}

func (repo *repository) GetSession(ctx context.Context, sessionID domain.UUID) (*domain.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)

	// Get the session.
	sessionJSON, err := repo.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrSessionNotFound
		}

		return nil, errors.Wrap(err, "cannot get")
	}

	// Decode the session.
	var session repoSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, errors.Wrap(err, "cannot decode")
	}

	return session.domainify(), nil
}
