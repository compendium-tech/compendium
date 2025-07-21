package model

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionLevel string

const (
	SubscriptionLevelStudent   SubscriptionLevel = "student"
	SubscriptionLevelTeam      SubscriptionLevel = "team"
	SubscriptionLevelCommunity SubscriptionLevel = "community"
)

type Subscription struct {
	UserID            uuid.UUID
	SubscriptionLevel SubscriptionLevel
	Till              time.Time
	Since             time.Time
}
