package model

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionLevel string

const (
	SubscriptionLevelFree    SubscriptionLevel = "student"
	SubscriptionLevelBasic   SubscriptionLevel = "team"
	SubscriptionLevelPremium SubscriptionLevel = "community"
)

type Subscription struct {
	UserID            uuid.UUID
	SubscriptionLevel SubscriptionLevel
	Till              time.Time
	Since             time.Time
}
