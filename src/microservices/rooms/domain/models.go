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
	Key string `json:"key" redis:",key"`
	Ver int64  `json:"ver" redis:",ver"`

	Name    string             `json:"name"`
	Viewers map[string]*Viewer `json:"viewers"`
}

type RoomInfo struct {
	ID      string             `json:"id,omitempty"`
	Name    string             `json:"name,omitempty"`
	Viewers map[string]*Viewer `json:"viewers,omitempty"`
}
