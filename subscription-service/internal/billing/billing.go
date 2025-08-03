package billing

import "context"

type Customer struct {
	ID    string
	Email string
}

type BillingAPI interface {
	GetCustomer(ctx context.Context, customerID string) (Customer, error)
	IsSubscriptionCanceled(ctx context.Context, subscriptionID string) (bool, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
}
