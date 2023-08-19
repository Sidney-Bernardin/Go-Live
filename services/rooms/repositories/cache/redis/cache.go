package redis

import (
	"context"
	"encoding/json"
	"rooms/configuration"
	"rooms/domain"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(config *configuration.Config) domain.CacheRepository {

	// Connect to Redis.
	client := redis.NewClient(&redis.Options{
		Addr:     config.CacheURL,
		Password: config.CachePassw,
	})

	return &cacheRepository{client}
}

// InsertRoom sets room.
func (repo *cacheRepository) InsertRoom(ctx context.Context, room *domain.Room) error {
	err := repo.client.HSet(ctx, "rooms:"+room.ID, room).Err()
	return errors.Wrap(err, "cannot set room")
}

// DeleteRoom deletes roomID's Room.
func (repo *cacheRepository) DeleteRoom(ctx context.Context, roomID string) error {
	err := repo.client.Del(ctx, "rooms:"+roomID).Err()
	return errors.Wrap(err, "cannot delete room")
}

// GetRoom gets roomID's Room.
func (repo *cacheRepository) GetRoom(ctx context.Context, roomID string) (room *domain.Room, err error) {
	err = repo.client.HGetAll(ctx, "rooms:"+roomID).Scan(&room)
	return room, errors.Wrap(err, "cannot get room")
}

// AddUserToRoom pushes userID to the Users list.
func (repo *cacheRepository) AddUserToRoom(ctx context.Context, roomID, userID string) error {
	err := repo.client.LPush(ctx, "rooms:"+roomID+":viewers", userID).Err()
	return errors.Wrap(err, "cannot push viewer")
}

// AddUserToRoom removes userID from the Users list.
func (repo *cacheRepository) RemoveUserFromRoom(ctx context.Context, roomID, userID string) error {
	err := repo.client.LRem(ctx, "rooms:"+roomID+":viewers", 0, userID).Err()
	return errors.Wrap(err, "cannot remove viewer")
}

func (repo *cacheRepository) SubToRoomEvents(ctx context.Context, eventChan chan domain.ChanMsg[*domain.RoomEvent], topic string) {

	sub := repo.client.Subscribe(ctx, topic)

	for {

		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			eventChan <- domain.ChanMsg[*domain.RoomEvent]{
				Err: errors.Wrap(err, "cannot receive message")}
			continue
		}

		var roomEvent *domain.RoomEvent
		if err := json.Unmarshal([]byte(msg.Payload), &roomEvent); err != nil {
			eventChan <- domain.ChanMsg[*domain.RoomEvent]{
				Err: errors.Wrap(err, "cannot decode room event")}
			continue
		}

		eventChan <- domain.ChanMsg[*domain.RoomEvent]{
			Content: roomEvent,
		}
	}
}

func (repo *cacheRepository) Publish(ctx context.Context, topic string, event any) error {

	b, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "cannot encode room event")
	}

	err = repo.client.Publish(ctx, topic, b).Err()
	return errors.Wrap(err, "cannot publish message")
}
