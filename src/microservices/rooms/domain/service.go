package domain

import (
	"sync"

	"github.com/pkg/errors"

	"rooms/configuration"
)

type Service interface {
	CreateRoom(sessionID, streamName string, settings *RoomSettings) error
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

func (svc *service) CreateRoom(sessionID, streamName string, settings *RoomSettings) error {

	user, err := svc.usersClientRepo.GetSelf(sessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	if user.ID != streamName {
		return ProblemDetail{
			Type:   PDTypeBadStreamID,
			Detail: "Your stream-name and user-ID dont't match.",
		}
	}

	room := &Room{
		Key:     user.ID,
		Name:    settings.Name,
		Viewers: map[string]*Viewer{},
	}

	if err := svc.cacheRepo.InsertRoom(room); err != nil {
		return errors.Wrap(err, "cannot insert room")
	}

	svc.mu.Lock()
	svc.roomChannels[user.ID] = make(chan map[string]any)
	svc.mu.Unlock()

	return nil
}

func (svc *service) GetRoom(roomID string) (*Room, error) {

	room, err := svc.cacheRepo.GetRoom(roomID)
	return room, errors.Wrap(err, "cannot get room")
}

func (svc *service) JoinRoom(sessionID, roomID string) (*User, <-chan map[string]any, error) {

	user, err := svc.usersClientRepo.GetSelf(sessionID, []string{"username"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get self")
	}

	viewer := &Viewer{
		UserID: user.ID,
	}

	if err := svc.cacheRepo.InsertViewer(roomID, viewer); err != nil {
		return nil, nil, errors.Wrap(err, "cannot insert viewer")
	}

	diagnosticsChan, ok := svc.roomChannels[roomID]
	if !ok {
		return nil, nil, ProblemDetail{
			Type:   PDTypeRoomDoesntExist,
			Detail: "The room's diagnostics weren't found.",
		}
	}

	return user, diagnosticsChan, nil
}

func (svc *service) LeaveRoom(userID, roomID string) error {

	err := svc.cacheRepo.DeleteViewer(userID, roomID)
	return errors.Wrap(err, "cannot delete viewer")
}

func (svc *service) DeleteRoom(sessionID string) error {

	user, err := svc.usersClientRepo.GetSelf(sessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	if err := svc.cacheRepo.DeleteRoom(user.ID); err != nil {
		return errors.Wrap(err, "cannot delete room")
	}

	svc.mu.Lock()
	delete(svc.roomChannels, user.ID)
	svc.mu.Unlock()

	return nil
}

func (svc *service) BroadcastMessage(user *User, roomID string, msg map[string]any) error {

	switch msg["type"] {

	case "CHAT":
		msg["user_id"] = user.ID
		msg["username"] = user.Username

	default:
		return nil
	}

	roomChan, ok := svc.roomChannels[roomID]
	if !ok {
		return ProblemDetail{
			Type: PDTypeRoomDoesntExist,
		}
	}

	room, err := svc.cacheRepo.GetRoom(roomID)
	if err != nil {
		return errors.Wrap(err, "cannot get room")
	}

	for range room.Viewers {
		roomChan <- msg
	}

	return nil
}
