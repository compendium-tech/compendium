package model

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a unique user login instance, tied to a specific device
// and browser/application.
//
// It's crucial for session management, allowing per-session rate limiting of
// refresh token requests and enabling users to manage (e.g., revoke) individual
// login sessions from a dashboard.
type Session struct {
	ID        uuid.UUID
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}

// RefreshToken represents a stored refresh token in the auth cache storage.
//
// ExpiresAt indicates the time when this refresh token should expire.
// It should not be used for validation after retrieval, as the storage
// mechanism handles expiration.
type RefreshToken struct {
	Session

	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}
