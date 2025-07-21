package domain

import (
	"time"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
)

type PutSubscriptionRequest struct {
	UserID            uuid.UUID
	SubscriptionLevel model.SubscriptionLevel
	Till              time.Time
	Since             time.Time
}
