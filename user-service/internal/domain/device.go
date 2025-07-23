package domain

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Os        string    `json:"os"`
	Device    string    `json:"device"`
	Location  string    `json:"location"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
