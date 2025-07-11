package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id              uuid.UUID
	Name            string
	Email           string
	IsEmailVerified bool
	IsAdmin         bool
	PasswordHash    []byte
	CreatedAt       time.Time
}
