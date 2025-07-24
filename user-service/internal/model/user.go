package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the system.
//
// It's crucial to understand the significance of the IsEmailVerified field.
// A User record where IsEmailVerified is false essentially represents an
// unconfirmed registration attempt. Such a user is considered to be in the
// sign-up process, awaiting email verification. Until the email is verified,
// this user record does not represent a fully registered or active user
// and should generally be ignored for most system operations.
//
// Furthermore, if a new sign-up attempt occurs with an email address
// that corresponds to an existing unverified user, the existing unverified
// record might be overwritten or updated (e.g., with a new password hash
// and creation timestamp) as part of issuing a new verification code.
type User struct {
	ID              uuid.UUID
	Name            string
	Email           string
	IsEmailVerified bool
	IsAdmin         bool
	PasswordHash    []byte
	CreatedAt       time.Time
}
