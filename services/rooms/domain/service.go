package domain

import (
	"sync"

	"github.com/pkg/errors"

	"rooms/configuration"
)

type Service interface {
	CreateRoom(streamName string, roomKey *RoomKey) error
	GetRoom(roomID string) (*Room, error)
	JoinRoom(sessionID, roomID string) (*User, <-chan map[string]any, error)
	LeaveRoom(userID, roomID string) error
	DeleteRoom(sessionID string) error

	BroadcastMessage(user *User, roomID string, msg map[string]any) error
}

type service struct {
	config *configuration.Configuration
	mu     *sync.Mutex

	cacheRepo       CacheRepository
	usersClientRepo UsersClientRepository

	roomChannels map[string]chan map[string]any
}

func NewService(
	config *configuration.Configuration,
	cacheRepo CacheRepository,
	usersClientRepo UsersClientRepository,
) Service {
	return &service{
		config:          config,
		mu:              &sync.Mutex{},
		cacheRepo:       cacheRepo,
		usersClientRepo: usersClientRepo,
		roomChannels:    make(map[string]chan map[string]any),
	}
}

func (svc *service) CreateRoom(streamName string, roomKey *RoomKey) error {

	user, err := svc.usersClientRepo.GetSelf(roomKey.SessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	// Make sure another user's stream-name isn't being used.
	if user.ID != streamName {
		return ProblemDetail{
			Type:   PDTypeBadStreamID,
			Detail: "Your stream-name and user-ID must match.",
		}
	}

	// Create a room.
	room := &Room{
		ID:      user.ID,
		Name:    roomKey.RoomSettings.Name,
		Viewers: map[string]*Viewer{},
	}

	// Insert the room into the cache.
	if err := svc.cacheRepo.InsertRoom(room); err != nil {
		return errors.Wrap(err, "cannot insert room")
	}

	// Create a channel for the room.
	svc.mu.Lock()
	svc.roomChannels[user.ID] = make(chan map[string]any)
	svc.mu.Unlock()

	return nil
}

func (svc *service) GetRoom(roomID string) (*Room, error) {

	// Get the room from the cache.
	room, err := svc.cacheRepo.GetRoom(roomID)
	return room, errors.Wrap(err, "cannot get room")
}

func (svc *service) JoinRoom(sessionID, roomID string) (*User, <-chan map[string]any, error) {

	user, err := svc.usersClientRepo.GetSelf(sessionID, []string{"username"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get self")
	}

	// Create a viewer.
	viewer := &Viewer{
		UserID: user.ID,
	}

	// Insert the viewer into the cache.
	if err := svc.cacheRepo.InsertViewer(roomID, viewer); err != nil {
		return nil, nil, errors.Wrap(err, "cannot insert viewer")
	}

	// Get the room's channel.
	roomChan, ok := svc.roomChannels[roomID]
	if !ok {
		return nil, nil, ProblemDetail{
			Type: PDTypeRoomDoesntExist,
		}
	}

	return user, roomChan, nil
}

func (svc *service) LeaveRoom(userID, roomID string) error {

	// Delete the viewer from the cache.
	err := svc.cacheRepo.DeleteViewer(userID, roomID)
	return errors.Wrap(err, "cannot delete viewer")
}

func (svc *service) DeleteRoom(sessionID string) error {

	user, err := svc.usersClientRepo.GetSelf(sessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	// Delete the room from the cache.
	if err := svc.cacheRepo.DeleteRoom(user.ID); err != nil {
		return errors.Wrap(err, "cannot delete room")
	}

	// Delete the room's channel.
	svc.mu.Lock()
	delete(svc.roomChannels, user.ID)
	svc.mu.Unlock()

	return nil
}

func (svc *service) BroadcastMessage(user *User, roomID string, msg map[string]any) error {

	switch msg["type"] {

	// For chat messages, add basic user fields to the message, to help other
	// users identify this user.
	case "CHAT":
		msg["user_id"] = user.ID
		msg["username"] = user.Username

	default:
		return nil
	}

	// Get the room's channel.
	roomChan, ok := svc.roomChannels[roomID]
	if !ok {
		return ProblemDetail{
			Type: PDTypeRoomDoesntExist,
		}
	}

	// Get the room from the cache.
	room, err := svc.cacheRepo.GetRoom(roomID)
	if err != nil {
		return errors.Wrap(err, "cannot get room")
	}

	// Send the message for each viewer in the room.
	for range room.Viewers {
		roomChan <- msg
	}

	return nil
}
