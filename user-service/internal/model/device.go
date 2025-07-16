package model

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	UserAgent string
	IpAddress string
	CreatedAt time.Time
}
