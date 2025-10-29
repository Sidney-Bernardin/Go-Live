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

func sessionKey(sessionID domain.UUID) string {
	return fmt.Sprintf("session:%s", sessionID)
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
	key := sessionKey(session.ID)
	p := repo.client.TxPipeline()

	// Set the session.
	p.JSONSet(ctx, key, ".", &repoSession{
		ID:        uuid.UUID(session.ID),
		UserID:    uuid.UUID(session.UserID),
		CSRFToken: string(session.CSRFToken),
	})

	// Set the session's expiry.
	p.Expire(ctx, key, repo.config.SessionDuration)

	// Execute the transaction.
	_, err := p.Exec(ctx)
	return errors.Wrap(err, "failed command executions")
}

func (repo *repository) GetSession(ctx context.Context, sessionID domain.UUID) (*domain.Session, error) {

	// Get the session.
	sessionJSON, err := repo.client.JSONGet(ctx, sessionKey(sessionID), ".").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrSessionNotFound
		}

		return nil, errors.Wrap(err, "failed getting")
	}

	// Decode the session.
	var session repoSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, errors.Wrap(err, "failed decoding")
	}

	return session.domainify(), nil
}
