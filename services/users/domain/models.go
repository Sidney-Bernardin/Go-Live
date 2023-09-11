package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      string             `json:"id,omitempty" bson:"-"`
	MongoID primitive.ObjectID `json:"-" bson:"_id"`

	ProfilePictureID      string             `json:"profile_picture_id,omitempty" bson:"-"`
	MongoProfilePictureID primitive.ObjectID `json:"-" bson:"profile_picture_id"`

	Username string `json:"username,omitempty" bson:"username"`
	Email    string `json:"email,omitempty" bson:"email"`
	Password string `json:"password,omitempty" bson:"password"`
}

type Session struct {
	ID      string             `json:"id,omitempty" bson:"-"`
	MongoID primitive.ObjectID `json:"-" bson:"_id"`

	UserID      string             `json:"user_id,omitempty" bson:"-"`
	MongoUserID primitive.ObjectID `json:"-" bson:"user_id"`

	ExpireAt time.Time `json:"expire_at,omitempty" bson:"expire_at"`
}

type SignupInfo struct {
	Username string `schema:"username,required"`
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

type SigninInfo struct {
	Username string `schema:"username,required"`
	Password string `schema:"password,required"`
}

type LoginResponse struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}
