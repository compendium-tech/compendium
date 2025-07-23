package domain

import (
	"time"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
)

type PutSubscriptionRequest struct {
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

type Subscription struct {
	Level model.SubscriptionLevel `json:"level"`
	Since time.Time               `json:"since"`
	Till  time.Time               `json:"till"`
}

type SubscriptionResponse struct {
	IsActive      bool `json:"isActive"`
	*Subscription `json:"subscription,omitempty"`
}
