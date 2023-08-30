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

// InsertRoom inserts room's basic fields as a redis hash.
func (repo *cacheRepository) InsertRoom(ctx context.Context, room *domain.Room) error {

	// Check if the Room already exists.
	exists, err := repo.client.Exists(ctx, "rooms:"+room.ID).Result()
	if err != nil {
		return errors.Wrap(err, "cannot check if room exists")
	}

	if exists == 1 {
		return domain.ProblemDetail{Problem: domain.ProblemRoomAlreadyExists}
	}

	// Inserts the Room's basic fields as a redis hash.
	err = repo.client.HSet(ctx, "rooms:"+room.ID, room).Err()
	return errors.Wrap(err, "cannot set room")
}

// DeleteRoom deletes roomID's Room.
func (repo *cacheRepository) DeleteRoom(ctx context.Context, roomID string) error {
	err := repo.client.Del(ctx, "rooms:"+roomID).Err()
	return errors.Wrap(err, "cannot delete room")
}

// GetRoom gets roomID's Room.
func (repo *cacheRepository) GetRoom(ctx context.Context, roomID string) (*domain.Room, error) {

	// Get the Room's hash.
	cmd := repo.client.HGetAll(ctx, "rooms:"+roomID)
	if err := cmd.Err(); err != nil {
		return nil, errors.Wrap(err, "cannot get room")
	}

	// Check if the Room wasn't found.
	if len(cmd.Val()) == 0 {
		return nil, domain.ProblemDetail{Problem: domain.ProblemRoomDoesntExist}
	}

	// Decode the Room's hash.
	var room domain.Room
	err := cmd.Scan(&room)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode room")
	}

	// Get the Room's list of viewers.
	room.Viewers, err = repo.client.LRange(ctx, "rooms:"+roomID+":viewers", 0, -1).Result()
	return &room, errors.Wrap(err, "cannot get list of viewers")
}

// AddViewerToRoom pushes userID to the viewers list of roomID's Room.
func (repo *cacheRepository) AddViewerToRoom(ctx context.Context, roomID, userID string) error {
	err := repo.client.RPush(ctx, "rooms:"+roomID+":viewers", userID).Err()
	return errors.Wrap(err, "cannot push viewer")
}

// AddViewerToRoom removes userID from the Users list.
func (repo *cacheRepository) RemoveViewerFromRoom(ctx context.Context, roomID, userID string) error {
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
