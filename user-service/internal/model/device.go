package model

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}
