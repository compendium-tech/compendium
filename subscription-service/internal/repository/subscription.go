package repository

import (
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	PutSubscription(sub model.Subscription) error
	GetSubscriptionByUserID(userID uuid.UUID) (*model.Subscription, error)
	RemoveSubscription(id string) error
}
