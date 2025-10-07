package domain

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func NewUUID() UUID {
	return UUID{uuid.New()}
}

func MustRandomString(length int) string {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
