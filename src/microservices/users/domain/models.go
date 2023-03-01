package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      string             `json:"id,omitempty" bson:"-"`
	MongoID primitive.ObjectID `json:"-" bson:"_id"`

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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
