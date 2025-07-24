package model

import (
	"time"

	"github.com/google/uuid"
)

// TrustedDevice represents a device that has been marked as trusted by a user.
// This model helps determine if Multi-Factor Authentication (MFA) is required
// for a login attempt from a specific device.
//
// It differs from the [Session] model (stored in refresh token repository) in its purpose:
// TrustedDevice is for MFA decisions, while [Session] is for user session management,
// allowing users to identify and revoke specific active sessions from their dashboard.
type TrustedDevice struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}
