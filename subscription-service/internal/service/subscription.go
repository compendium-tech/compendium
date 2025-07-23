package service

import (
	"context"

	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
)

type SubscriptionService interface {
	PutSubscription(ctx context.Context, request domain.PutSubscriptionRequest) error
	CancelSubscription(ctx context.Context, subscriptionID string) error
}

type subscriptionService struct {
	billingAPI             billing.BillingAPI
	productIDs             config.ProductIDs
	userService            interop.UserService
	billingLockRepository  repository.BillingLockRepository
	subscriptionRepository repository.SubscriptionRepository
}

func NewSubscriptionService(
	billingAPI billing.BillingAPI,
	productIDs config.ProductIDs,
	userService interop.UserService,
	billingLockRepository repository.BillingLockRepository,
	subscriptionRepository repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		billingAPI:             billingAPI,
		productIDs:             productIDs,
		userService:            userService,
		billingLockRepository:  billingLockRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *subscriptionService) PutSubscription(ctx context.Context, request domain.PutSubscriptionRequest) (finalErr error) {
	if len(request.Items) == 0 {
		return appErr.Errorf(appErr.RequestValidationError, "No items in request")
	}

	customer, err := s.billingAPI.GetCustomer(ctx, request.CustomerID)
	if err != nil {
		return err
	}

	lock, err := s.billingLockRepository.ObtainLock(ctx, request.CustomerID)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	account, err := s.userService.FindAccountByEmail(ctx, customer.Email)
	if err != nil {
		return err
	}

	if account == nil {
		return appErr.Errorf(appErr.UserNotFoundError, "Account not found for email: %s", customer.Email)
	}

	if len(request.Items) > 1 {
		log.L(ctx).Warnf("Multiple items in request, cancelling subscription for user %s", account.ID)

		err := s.billingAPI.CancelSubscription(ctx, request.SubscriptionID)
		if err != nil {
			return err
		}
	}

	item := request.Items[0]

	subscription, err := s.subscriptionRepository.GetSubscriptionByUserID(account.ID)
	if err != nil {
		return err
	}

	if subscription != nil {
		err := s.billingAPI.CancelSubscription(ctx, subscription.ID)

		if err != nil {
			log.L(ctx).Errorf("Failed to cancel existing subscription for user %s: %v", account.ID, err)

			return err
		}
	}

	var subscriptionLevel model.SubscriptionLevel

	switch item.ProductID {
	case s.productIDs.StudentSubscriptionProductID:
		subscriptionLevel = model.SubscriptionLevelStudent
	case s.productIDs.TeamSubscriptionProductID:
		subscriptionLevel = model.SubscriptionLevelTeam
	case s.productIDs.CommunitySubscriptionProductID:
		subscriptionLevel = model.SubscriptionLevelCommunity
	default:
		{
			log.L(ctx).Errorf("Unknown product ID %s for user %s, cancelling subscription", item.ProductID, account.ID)

			err := s.billingAPI.CancelSubscription(ctx, request.SubscriptionID)
			if err != nil {
				return err
			}
		}
	}

	return s.subscriptionRepository.PutSubscription(model.Subscription{
		ID:                request.SubscriptionID,
		UserID:            account.ID,
		SubscriptionLevel: subscriptionLevel,
		Till:              request.Till,
		Since:             request.Since,
	})
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, subscriptionID string) error {
	return s.subscriptionRepository.RemoveSubscription(subscriptionID)
}
