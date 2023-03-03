package domain

import (
	"sync"

	"github.com/pkg/errors"

	"rooms/configuration"
)

type Service interface {
	CreateRoom(sessionID string, settings *RoomSettings) error
	GetRoom(roomID string) (*Room, error)
	JoinRoom(sessionID, roomID string) error
	LeaveRoom(sessionID, roomID string) error
	DeleteRoom(sessionID string) error

	NewDiagnosticsListener(sessionID, roomID string) (*User, <-chan *Diagnostic, error)
	BroadcastDiagnostic(user *User, roomID string, diagnostic *Diagnostic) error
}

type service struct {
	config *configuration.Configuration

	mu *sync.Mutex

	cacheRepo       CacheRepository
	usersClientRepo UsersClientRepository

	diagnosticsChans map[string]chan *Diagnostic
}

// NewService returns a new Service.
func NewService(
	config *configuration.Configuration,
	cacheRepo CacheRepository,
	usersClientRepo UsersClientRepository,
) Service {
	return &service{
		config:           config,
		mu:               &sync.Mutex{},
		cacheRepo:        cacheRepo,
		usersClientRepo:  usersClientRepo,
		diagnosticsChans: map[string]chan *Diagnostic{},
	}
}

func (svc *service) CreateRoom(sessionID string, settings *RoomSettings) error {

	user, err := svc.usersClientRepo.GetSelf(sessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	// Create a room.
	room := &Room{
		Key:     user.ID,
		Name:    settings.Name,
		Viewers: map[string]*Viewer{},
	}

	// Insert the room into the cache.
	if err := svc.cacheRepo.InsertRoom(room); err != nil {
		return errors.Wrap(err, "cannot insert room")
	}

	// Create a diagnostics channel for the room.
	svc.mu.Lock()
	svc.diagnosticsChans[user.ID] = make(chan *Diagnostic)
	svc.mu.Unlock()
	return nil
}

func (svc *service) GetRoom(roomID string) (*Room, error) {

	// Get the room from the cache.
	room, err := svc.cacheRepo.GetRoom(roomID)
	return room, errors.Wrap(err, "cannot get room")
}

func (svc *service) JoinRoom(sessionID, roomID string) error {

	user, err := svc.usersClientRepo.GetSelf(sessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	// Create a viewer and insert it into the cached room.
	err = svc.cacheRepo.InsertViewer(roomID, &Viewer{user.ID})
	return errors.Wrap(err, "cannot insert viewer")
}

func (svc *service) LeaveRoom(sessionID, roomID string) error {

	user, err := svc.usersClientRepo.GetSelf(sessionID, nil)
	if err != nil {
		return errors.Wrap(err, "cannot get self")
	}

	// Delete the user's viewer from the cached room.
	err = svc.cacheRepo.DeleteViewer(roomID, user.ID)
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

	// Delete the room's diagnostics channel.
	svc.mu.Lock()
	delete(svc.diagnosticsChans, user.ID)
	svc.mu.Unlock()

	return nil
}

func (svc *service) NewDiagnosticsListener(sessionID, roomID string) (*User, <-chan *Diagnostic, error) {

	user, err := svc.usersClientRepo.GetSelf(sessionID, []string{"username"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get self")
	}

	// Make sure the user is in the room by getting the user's viewer.
	if _, err := svc.cacheRepo.GetViewer(roomID, user.ID); err != nil {
		return nil, nil, errors.Wrap(err, "cannot get viewer")
	}

	// Get the room's diagnostics channel.
	diagnosticsChan, ok := svc.diagnosticsChans[roomID]
	if !ok {
		return nil, nil, ProblemDetail{
			Type:   PDTypeRoomDoesntExist,
			Detail: "The room's diagnostics weren't found.",
		}
	}

	return user, diagnosticsChan, nil
}

func (svc *service) BroadcastDiagnostic(user *User, roomID string, diagnostic *Diagnostic) error {

	switch diagnostic.Type {

	case "CHAT_MESSAGE":
		diagnostic.ChatMessage.UserID = user.ID
		diagnostic.ChatMessage.Username = user.Username

	default:
		return nil
	}

	// Get the room's diagnostics channel.
	diagnosticsChan, ok := svc.diagnosticsChans[roomID]
	if !ok {
		return ProblemDetail{
			Type:   PDTypeRoomDoesntExist,
			Detail: "The room's diagnostics weren't found.",
		}
	}

	// Get the room from the cache.
	room, err := svc.cacheRepo.GetRoom(roomID)
	if err != nil {
		return errors.Wrap(err, "cannot get room")
	}

	// Broadcast the diagnostic by sending it through the diagnostics channel
	// for each viewer in the room.
	for range room.Viewers {
		diagnosticsChan <- diagnostic
	}

	return nil
}
