package repository

import (
	"github.com/adslmgrv/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	CreateSubscription(sub model.Subscription) error
	GetSubscriptionByUserID(userID uuid.UUID) (*model.Subscription, error)
	RemoveSubscription(userID uuid.UUID) error
}
