package domain

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Room struct {
	ID      string     `json:"id,omitempty" redis:"key1,id"`
	Name    string     `json:"name,omitempty" redis:"name,key2"`
	Viewers []string`json:"users,omitempty"`
}

type ChanMsg[T any] struct {
	Content T     `json:"content,omitempty"`
	Err     error `json:"error,omitempty"`
}

type RoomEvent struct {
	Type string `json:"type"`
	ChatMsg
}

type ChatMsg struct {
	UserID  string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
}
