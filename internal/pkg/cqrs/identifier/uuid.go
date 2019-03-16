package identifier

import "github.com/google/uuid"

type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func (u UUID) String() string {
	return u.String()
}
