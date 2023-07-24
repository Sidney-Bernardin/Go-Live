package domain

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Viewer struct {
	UserID string `json:"user_id"`
}

type Room struct {
	ID  string `json:"id,omitempty"`
	Key string `json:"key" redis:",key"`
	Ver int64  `json:"ver" redis:",ver"`

	Name    string             `json:"name"`
	Viewers map[string]*Viewer `json:"viewers"`
}

type RoomKey struct {
	SessionID    string `json:"session_id"`
	RoomSettings struct {
		Name string `json:"name"`
	} `json:"room_settings"`
}
