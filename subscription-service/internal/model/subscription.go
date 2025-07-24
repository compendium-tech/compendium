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

func (s SubscriptionLevel) Priority() int {
	switch s {
	case SubscriptionLevelStudent:
		return 1
	case SubscriptionLevelTeam:
		return 2
	case SubscriptionLevelCommunity:
		return 3
	default:
		return 0
	}
}

type Subscription struct {
	ID             string
	BackedBy       uuid.UUID
	Level          SubscriptionLevel
	InvitationCode *string
	Till           time.Time
	Since          time.Time
}

type SubscriptionMember struct {
	SubscriptionID string
	UserID         uuid.UUID
	Since          time.Time
}
