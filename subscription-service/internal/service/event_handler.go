package service

import (
	"context"
	"fmt"

	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
)

type BillingEventHandlerService interface {
	HandleUpdatedSubscription(ctx context.Context, request domain.HandleUpdatedSubscriptionRequest)
	CancelSubscription(ctx context.Context, subscriptionID string)
}

type billingEventHandlerService struct {
	billingAPI             billing.BillingAPI
	productIDs             config.ProductIDs
	billingLockRepository  repository.BillingLockRepository
	subscriptionRepository repository.SubscriptionRepository
}

func NewBillingEventHandlerService(
	billingAPI billing.BillingAPI,
	productIDs config.ProductIDs,
	billingLockRepository repository.BillingLockRepository,
	subscriptionRepository repository.SubscriptionRepository) BillingEventHandlerService {
	return &billingEventHandlerService{
		billingAPI:             billingAPI,
		productIDs:             productIDs,
		billingLockRepository:  billingLockRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *billingEventHandlerService) HandleUpdatedSubscription(ctx context.Context, request domain.HandleUpdatedSubscriptionRequest) {
	logger := log.L(ctx).WithField("subscriptionId", request.SubscriptionID).WithField("userId", request.UserID)
	logger.Info("Upserting subscription")

	if len(request.Items) == 0 {
		logger.Warn("No items in upsert subscription request")
		myerror.NewWithReason(myerror.RequestValidationError, "No items in request").Throw()
	}

	lock := s.billingLockRepository.ObtainLock(ctx, request.UserID)
	defer lock.Release(ctx)

	for i, item := range request.Items {
		itemLogger := logger.WithField("itemIndex", i).WithField("productID", item.ProductID)
		itemLogger.Info("Processing subscription item")

		var subscriptionLevel model.Tier

		switch item.ProductID {
		case s.productIDs.StudentSubscriptionProductID:
			subscriptionLevel = model.TierStudent
		case s.productIDs.TeamSubscriptionProductID:
			subscriptionLevel = model.TierTeam
		case s.productIDs.CommunitySubscriptionProductID:
			subscriptionLevel = model.TierCommunity
		default:
			itemLogger.Errorf("Unknown product ID %s, skipping this item", item.ProductID)

			if len(request.Items) == 1 {
				if !s.billingAPI.IsSubscriptionCanceled(ctx, request.SubscriptionID) {
					s.billingAPI.CancelSubscription(ctx, request.SubscriptionID)
				}

				myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("Unknown product ID: %s", item.ProductID)).Throw()
			}
			continue
		}

		// To maintain a single active subscription per payer within this service,
		// we first cancel any existing subscription before upserting a new one.
		//
		// If the cancellation fails, the external billing service will retry the request,
		// ensuring the database isn't updated with a new subscription until the previous one is successfully canceled.
		existingSubscription := s.subscriptionRepository.FindSubscriptionByPayerUserID(ctx, request.UserID)

		if existingSubscription != nil {
			itemLogger.Infof("Existing subscription %s found, cancelling it before upserting new item", existingSubscription.ID)

			if !s.billingAPI.IsSubscriptionCanceled(ctx, existingSubscription.ID) {
				s.billingAPI.CancelSubscription(ctx, existingSubscription.ID)
			}
		}

		s.subscriptionRepository.UpsertSubscription(ctx, model.Subscription{
			ID:       request.SubscriptionID,
			BackedBy: request.UserID,
			Tier:     subscriptionLevel,
			Till:     request.Till,
			Since:    request.Since,
		})

		itemLogger.Info("Subscription item processed successfully")
	}

	logger.Info("All subscription items processed, upsert completed successfully")
}

func (s *billingEventHandlerService) CancelSubscription(ctx context.Context, subscriptionID string) {
	logger := log.L(ctx).WithField("subscriptionId", subscriptionID)
	logger.Info("Cancelling subscription")

	s.subscriptionRepository.RemoveSubscription(ctx, subscriptionID)

	log.L(ctx).Info("Subscription cancellation initiated successfully")
}
