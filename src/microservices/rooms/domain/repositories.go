package domain

type CacheRepository interface {
	InsertRoom(room *Room) error
	GetRoom(roomID string) (*RoomInfo, error)
	DeleteRoom(roomID string) error

	InsertViewer(roomID string, viewer *Viewer) error
	GetViewer(roomID, viewerID string) (*Viewer, error)
	DeleteViewer(roomID, viewerID string) error
}

type UsersClientRepository interface {
	GetSelf(sessionID string, fields []string) (*User, error)
}
