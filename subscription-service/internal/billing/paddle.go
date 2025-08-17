package billing

import (
	"context"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/compendium-tech/compendium/common/pkg/log"
)

type paddleBilling struct {
	sdk paddle.SDK
}

func NewPaddleBillingAPI(sdk paddle.SDK) BillingAPI {
	return &paddleBilling{
		sdk: sdk,
	}
}

func (pb *paddleBilling) GetCustomer(ctx context.Context, customerID string) Customer {
	customer, err := pb.sdk.GetCustomer(ctx, &paddle.GetCustomerRequest{
		CustomerID: customerID,
	})
	if err != nil {
		panic(err)
	}

	return Customer{
		ID:    customer.ID,
		Email: customer.Email,
	}
}

func (pb *paddleBilling) IsSubscriptionCanceled(ctx context.Context, subscriptionID string) bool {
	subscription, err := pb.sdk.GetSubscription(ctx, &paddle.GetSubscriptionRequest{
		SubscriptionID: subscriptionID,
	})
	if err != nil {
		log.L(ctx).
			WithField("subscriptionId", subscriptionID).
			Warnf("Could not get subscription information: %s", err.Error())
		return false
	}

	if subscription == nil {
		return true
	}

	if subscription.Status == paddle.SubscriptionStatusCanceled {
		return true
	}

	return false
}

func (pb *paddleBilling) CancelSubscription(ctx context.Context, subscriptionID string) {
	immediately := paddle.EffectiveFromImmediately
	_, err := pb.sdk.CancelSubscription(ctx, &paddle.CancelSubscriptionRequest{
		SubscriptionID: subscriptionID,
		EffectiveFrom:  &immediately,
	})
	if err != nil {
		panic(err)
	}
}
