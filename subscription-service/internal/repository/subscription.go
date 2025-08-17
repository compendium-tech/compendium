package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
)

// SubscriptionRepository defines the interface for data access operations related to subscriptions and their members.
//
// UpsertSubscription creates a new subscription if one with the given ID does not exist, or updates an existing
// subscription if it does.
type SubscriptionRepository interface {
	UpsertSubscription(ctx context.Context, subscription model.Subscription)

	FindSubscriptionByInvitationCode(ctx context.Context, invitationCode string) *model.Subscription
	FindSubscriptionByMemberUserID(ctx context.Context, userID uuid.UUID) *model.Subscription
	FindSubscriptionByPayerUserID(ctx context.Context, backedBy uuid.UUID) *model.Subscription

	GetSubscriptionMembers(ctx context.Context, subscriptionID string) []model.SubscriptionMember

	CreateSubscriptionMemberAndCheckMemberCount(ctx context.Context, member model.SubscriptionMember, checkCount func(uint) error)
	RemoveSubscriptionMemberBySubscriptionAndUserID(ctx context.Context, subscriptionID string, userID uuid.UUID)

	RemoveSubscription(ctx context.Context, id string)
}
