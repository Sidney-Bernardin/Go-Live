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

type RoomSettings struct {
	Name string `json:"name"`
}

type Diagnostic struct {
	Type string `json:"type"`

	ChatMessage *struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Text     string `json:"text"`
	} `json:"chat_message,omitempty"`
}
