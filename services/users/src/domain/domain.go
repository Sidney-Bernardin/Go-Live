package domain

import (
	googleUUID "github.com/google/uuid"
)

type UUID googleUUID.UUID

func NewUUID() UUID {
	return UUID(googleUUID.New())
}

func NewUUIDFromString(str string) (UUID, error) {
	uuid, err := googleUUID.Parse(str)
	if err != nil {
		return UUID{}, &DomainError{DomainErrorTypeUUIDInvalid, err.Error(), nil}
	}
	return UUID(uuid), nil
}

func (uuid UUID) String() string {
	return googleUUID.UUID(uuid).String()
}
