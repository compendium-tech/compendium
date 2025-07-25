package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	IsCurrent bool      `json:"isCurrent"`
	Name      string    `json:"name"`
	Os        string    `json:"os"`
	Device    string    `json:"device"`
	Location  string    `json:"location"`
	UserAgent string    `json:"userAgent"`
	IPAddress string    `json:"ipAddress"`
	CreatedAt time.Time `json:"createdAt"`
}
