package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserId uuid.UUID
	Token  string
	Expiry time.Time
}
