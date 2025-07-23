package repository

import (
	"context"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	PutSubscription(ctx context.Context, sub model.Subscription) error
	GetSubscriptionByUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error)
	RemoveSubscription(ctx context.Context, id string) error
}
