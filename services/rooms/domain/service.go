package domain

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"rooms/configuration"
)

type Service interface {
	CreateRoom(ctx context.Context, sessionID, streamID string) error
	DeleteRoom(ctx context.Context, sessionID string) error
	GetRoom(ctx context.Context, roomID string) (*Room, error)

	JoinRoom(ctx context.Context, sessionID, roomID string) (userID string, err error)
	LeaveRoom(ctx context.Context, userID, roomID string) error

	ListenForRoomEvents(ctx context.Context, eventChan chan ChanMsg[*RoomEvent], roomID string)
	SendRoomEvent(ctx context.Context, userID, roomID string, event *RoomEvent) error
}

type service struct {
	config *configuration.Config

	cacheRepo       CacheRepository
	usersClientRepo UsersClientRepository
}

func NewService(
	config *configuration.Config,
	cacheRepo CacheRepository,
	usersClientRepo UsersClientRepository,
) Service {
	return &service{
		config:          config,
		cacheRepo:       cacheRepo,
		usersClientRepo: usersClientRepo,
	}
}

// CreateRoom creates a new Room for sessionID's User and inserts it into the cache.
func (svc *service) CreateRoom(ctx context.Context, sessionID, videoStreamID string) error {

	user, err := svc.usersClientRepo.AuthenticateUser(ctx, sessionID, "username")
	if err != nil {
		return errors.Wrap(err, "cannot get authenticate user")
	}

	// Make sure the User doesn't create a room for someone elses account.
	if user.ID != videoStreamID {
		return ProblemDetail{
			Problem: ProblemUnauthorized,
		}
	}

	// Insert a new Room for the User into the cache.
	err = svc.cacheRepo.InsertRoom(ctx, &Room{
		ID:   user.ID,
		Name: fmt.Sprintf("%s's Room", user.Username),
	})

	return errors.Wrap(err, "cannot insert room")
}

// DeleteRoom deletes the Room of sessionID's User from the cache.
func (svc *service) DeleteRoom(ctx context.Context, sessionID string) error {

	user, err := svc.usersClientRepo.AuthenticateUser(ctx, sessionID)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	// Delete the User-ID's Room from the cache.
	err = svc.cacheRepo.DeleteRoom(ctx, user.ID)
	return errors.Wrap(err, "cannot delete room")
}

// GetRoom roomID's Room from the cache.
func (svc *service) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	room, err := svc.cacheRepo.GetRoom(ctx, roomID)
	return room, errors.Wrap(err, "cannot get room")
}

// JoinRoom adds the User-ID of sessionID's User to the viewers list of
// roomID's Room in the cache.
func (svc *service) JoinRoom(ctx context.Context, sessionID, roomID string) (string, error) {

	user, err := svc.usersClientRepo.AuthenticateUser(ctx, sessionID)
	if err != nil {
		return "", errors.Wrap(err, "cannot get self")
	}

	// Add the User-ID to the Room in the cache.
	err = svc.cacheRepo.AddViewerToRoom(ctx, roomID, user.ID)
	return user.ID, errors.Wrap(err, "cannot insert viewer")
}

// LeaveRoom removes userID from the viewers list of roomID's Room in the cache.
func (svc *service) LeaveRoom(ctx context.Context, userID, roomID string) error {
	err := svc.cacheRepo.RemoveViewerFromRoom(ctx, roomID, userID)
	return errors.Wrap(err, "cannot delete viewer")
}

// ListenForRoomEvents sends the RoomEvents of the roomID's Room to eventChan.
func (svc *service) ListenForRoomEvents(ctx context.Context, eventChan chan ChanMsg[*RoomEvent], roomID string) {
	svc.cacheRepo.SubToRoomEvents(ctx, eventChan, roomID)
}

// SendRoomEvent sends event to the topic for roomID's Room in the cache.
func (svc *service) SendRoomEvent(ctx context.Context, userID, roomID string, event *RoomEvent) error {
	err := svc.cacheRepo.Publish(ctx, roomID, event)
	return errors.Wrap(err, "cannot publish event")
}
