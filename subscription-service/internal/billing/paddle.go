package billing

import (
	"context"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
)

type paddleBilling struct {
	sdk paddle.SDK
}

func NewPaddleBillingAPI(sdk paddle.SDK) BillingAPI {
	return &paddleBilling{
		sdk: sdk,
	}
}

func (pb *paddleBilling) GetCustomer(ctx context.Context, customerID string) (Customer, error) {
	customer, err := pb.sdk.GetCustomer(ctx, &paddle.GetCustomerRequest{
		CustomerID: customerID,
	})
	if err != nil {
		return Customer{}, err
	}

	return Customer{
		ID:    customer.ID,
		Email: customer.Email,
	}, nil
}

func (pb *paddleBilling) CancelSubscription(ctx context.Context, subscriptionID string) error {
	_, err := pb.sdk.CancelSubscription(ctx, &paddle.CancelSubscriptionRequest{
		SubscriptionID: subscriptionID,
	})
	if err != nil {
		return err
	}

	return nil
}
