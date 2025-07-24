package repository

import (
	"context"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	PutSubscription(ctx context.Context, sub model.Subscription) error
	GetSubscriptionByInvitationCode(ctx context.Context, invitationCode string) (*model.Subscription, error)
	GetSubscriptionByMemberUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error)
	GetSubscriptionByPayerUserID(ctx context.Context, backedBy uuid.UUID) (*model.Subscription, error)
	GetSubscriptionMembers(ctx context.Context, subscriptionID string) ([]model.SubscriptionMember, error)
	AddSubscriptionMember(ctx context.Context, member model.SubscriptionMember) error
	UpdateSubscriptionInvitationCode(ctx context.Context, subscriptionID string, code *string) error
	RemoveSubscriptionMemberBySubscriptionAndUserID(ctx context.Context, subscriptionID string, userID uuid.UUID) error
	RemoveSubscription(ctx context.Context, id string) error
}
