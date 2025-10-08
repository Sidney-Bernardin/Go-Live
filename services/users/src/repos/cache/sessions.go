package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type cacheSession struct {
	ID        domain.UUID `json:"id"`
	UserID    domain.UUID `json:"user_id"`
	CSRFToken string      `json:"csrf_token"`
}

func (s *cacheSession) domainify() *domain.Session {
	if s == nil {
		return nil
	}

	return &domain.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		CSRFToken: domain.CSRFToken(s.CSRFToken),
	}
}

func (c *cache) InsertSession(ctx context.Context, session *domain.Session) error {
	key := fmt.Sprintf("session:%s", session.ID)
	p := c.client.TxPipeline()

	jsonSetCmd := p.JSONSet(ctx, key, ".", &cacheSession{
		ID:        session.ID,
		UserID:    session.UserID,
		CSRFToken: string(session.CSRFToken),
	})
	expireCmd := p.Expire(ctx, key, c.config.SessionDuration)

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

func (c *cache) GetSession(ctx context.Context, sessionID domain.UUID) (*domain.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)

	sessionJSON, err := c.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrSessionNotFound
		}

		return nil, errors.Wrap(err, "cannot get")
	}

	var session cacheSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, errors.Wrap(err, "cannot decode")
	}

	return session.domainify(), nil
}
