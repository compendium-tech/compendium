package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    uuid.UUID
	Token     string
	SessionID uuid.UUID
	Expiry    time.Time
}
