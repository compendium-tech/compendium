package billing

import "context"

type Customer struct {
	ID    string
	Email string
}

type BillingAPI interface {
	GetCustomer(ctx context.Context, customerID string) (Customer, error)
	CancelSubscription(ctx context.Context, subscriptionId string) error
}
