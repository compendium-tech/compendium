package domain

import (
	"time"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
)

const (
	SubscriptionRolePayer  SubscriptionRole = "payer"
	SubscriptionRoleMember SubscriptionRole = "member"
)

type HandleUpdatedSubscriptionRequest struct {
	SubscriptionID string
	CustomerID     string
	Items          []SubscriptionItem
	Till           time.Time
	Since          time.Time
}

type SubscriptionItem struct {
	PriceID   string
	ProductID string
	Quantity  int
}

type SubscriptionRole string

type Subscription struct {
	Role    SubscriptionRole     `json:"role"`
	Level   model.Tier           `json:"level"`
	Since   time.Time            `json:"since"`
	Till    time.Time            `json:"till"`
	Members []SubscriptionMember `json:"members,omitempty"`
}

type SubscriptionResponse struct {
	IsActive      bool `json:"isActive"`
	*Subscription `json:"subscription,omitempty"`
}

type InvitationCodeResponse struct {
	InvitationCode *string `json:"invitationCode,omitempty"`
}

type SubscriptionMember struct {
	UserID          uuid.UUID        `json:"userId"`
	Name            string           `json:"name"`
	Email           string           `json:"email,omitempty"`
	Role            SubscriptionRole `json:"role,omitempty"`
	IsAccountActive bool             `json:"isAccountActive,omitempty"`
}
