package billing

import (
	"context"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/ztrue/tracerr"
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
	subscription, err := pb.sdk.GetSubscription(ctx, &paddle.GetSubscriptionRequest{
		SubscriptionID: subscriptionID,
	})

	if subscription == nil && err == nil {
		return nil
	}

	if subscription.Status == paddle.SubscriptionStatusCanceled {
		return nil
	}

	immediately := paddle.EffectiveFromImmediately
	_, err = pb.sdk.CancelSubscription(ctx, &paddle.CancelSubscriptionRequest{
		SubscriptionID: subscriptionID,
		EffectiveFrom:  &immediately,
	})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
