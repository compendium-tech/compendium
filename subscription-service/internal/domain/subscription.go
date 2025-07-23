package domain

import (
	"time"
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
