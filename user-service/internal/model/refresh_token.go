package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}

type RefreshToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	Session   Session
}
