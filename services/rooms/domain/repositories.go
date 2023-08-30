package domain

import "context"

type UsersClientRepository interface {
	AuthenticateUser(ctx context.Context, sessionID string, fields ...string) (*User, error)
}

type CacheRepository interface {
	InsertRoom(ctx context.Context, room *Room) error
	GetRoom(ctx context.Context, roomID string) (*Room, error)
	DeleteRoom(ctx context.Context, roomID string) error

	AddViewerToRoom(ctx context.Context, roomID, userID string) error
	RemoveViewerFromRoom(ctx context.Context, roomID, userID string) error

	SubToRoomEvents(ctx context.Context, eventChan chan ChanMsg[*RoomEvent], roomID string)
	Publish(ctx context.Context, topic string, event any) error
}
