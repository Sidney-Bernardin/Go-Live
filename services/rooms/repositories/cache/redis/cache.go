package redis

import (
	"context"
	"rooms/configuration"
	"rooms/domain"

	"github.com/pkg/errors"
	"github.com/rueian/rueidis"
	"github.com/rueian/rueidis/om"
)

type cacheRepository struct {
	config    *configuration.Configuration
	client    rueidis.Client
	roomsRepo om.Repository[domain.Room]
}

func NewCacheRepository(config *configuration.Configuration) (domain.CacheRepository, error) {

	// Connect to Redis.
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{config.RedisURL},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to redis")
	}

	// Create an object-map for rooms.
	roomsRepo := om.NewJSONRepository("rooms", domain.Room{}, client)

	// err = roomsRepo.CreateIndex(context.Background(), func(schema om.FtCreateSchema) om.Completed {
	// 	return schema.FieldName("$.name").Text().Build()
	// })

	// if err != nil {
	// 	return nil, errors.Wrap(err, "cannot create index for rooms")
	// }

	return &cacheRepository{config, client, roomsRepo}, nil
}

func (repo *cacheRepository) InsertRoom(room *domain.Room) error {

	// Create a copy of the room with the rooms object-map.
	newRoom := repo.roomsRepo.NewEntity()
	newRoom.Key = room.ID
	newRoom.Name = room.Name
	newRoom.Viewers = room.Viewers

	// Set the copy.
	if err := repo.roomsRepo.Save(context.Background(), newRoom); err != nil {

		// Check if the room already exists.
		if err.Error() == "object version mismatched, please retry" {
			return domain.ProblemDetail{
				Type: domain.PDTypeRoomAlreadyExists,
			}
		}

		return errors.Wrap(err, "cannot save room")
	}

	return nil
}

func (repo *cacheRepository) GetRoom(roomID string) (*domain.Room, error) {

	// Get the room.
	room, err := repo.roomsRepo.Fetch(context.Background(), roomID)
	if err != nil {

		// Check if the room wasn't found.
		if rueidis.IsRedisNil(err) {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeRoomDoesntExist,
			}
		}

		return nil, errors.Wrap(err, "cannot fetch room")
	}

	return &domain.Room{
		ID:      room.Key,
		Name:    room.Name,
		Viewers: room.Viewers,
	}, nil
}

func (repo *cacheRepository) DeleteRoom(roomID string) error {

	// Delete the room.
	err := repo.roomsRepo.Remove(context.Background(), roomID)
	return errors.Wrap(err, "cannot remove room")
}

func (repo *cacheRepository) InsertViewer(roomID string, viewer *domain.Viewer) error {

	// Sets the viewer into the room.
	cmd := repo.client.B().JsonSet().
		Key("rooms:" + roomID).
		Path("$.viewers." + viewer.UserID).
		Value(rueidis.JSON(viewer)).
		Nx().Build()

	// Use the command.
	if err := repo.client.Do(context.Background(), cmd).Error(); err != nil {

		// Check if the room wasn't found.
		if err.Error() == "new objects must be created at the root" {
			return domain.ProblemDetail{
				Type: domain.PDTypeViewerDoesntExist,
			}
		}

		// Check if the viewer already exists.
		if rueidis.IsRedisNil(err) {
			return domain.ProblemDetail{
				Type: domain.PDTypeViewerAlreadyExists,
			}
		}

		return errors.Wrap(err, "cannot set viewer into room")
	}

	return nil
}

func (repo *cacheRepository) GetViewer(roomID, userID string) (*domain.Viewer, error) {

	// Gets the viewer from the room.
	cmd := repo.client.B().JsonGet().
		Key("rooms:" + roomID).
		Paths("$.viewers." + userID).Build()

	// Use the command.
	var viewer domain.Viewer
	if err := repo.client.Do(context.Background(), cmd).DecodeJSON(&viewer); err != nil {

		// Check if the room wasn't found.
		if err.Error() == "new objects must be created at the root" {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeRoomDoesntExist,
			}
		}

		// Check if the viewer wasn't found.
		if rueidis.IsRedisNil(err) {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeViewerDoesntExist,
			}
		}

		return &viewer, errors.Wrap(err, "cannot get viewer")
	}

	return &viewer, nil
}

func (repo *cacheRepository) DeleteViewer(roomID, userID string) error {

	// Deletes the viewer.
	cmd := repo.client.B().JsonDel().
		Key("rooms:" + roomID).
		Path("$.viewers." + userID).Build()

	// Use the command.
	if err := repo.client.Do(context.Background(), cmd).Error(); err != nil {

		// Check if the room wasn't found.
		if err.Error() == "new objects must be created at the root" {
			return domain.ProblemDetail{
				Type: domain.PDTypeRoomDoesntExist,
			}
		}

		// Check if the viewer wasn't found.
		if rueidis.IsRedisNil(err) {
			return domain.ProblemDetail{
				Type: domain.PDTypeViewerDoesntExist,
			}
		}

		return errors.Wrap(err, "cannot delete viewer")
	}

	return nil
}
