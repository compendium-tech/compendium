package service

import (
	"context"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetSubscription(ctx context.Context) (*domain.SubscriptionResponse, error)
	CancelSubscription(ctx context.Context) error
	RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID) error

	PutSubscription(ctx context.Context, request domain.PutSubscriptionRequest) error
	RemoveSubscription(ctx context.Context, subscriptionID string) error
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

func (s *subscriptionService) GetSubscription(ctx context.Context) (*domain.SubscriptionResponse, error) {
	userId, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	subscription, err := s.subscriptionRepository.GetSubscriptionByMemberUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return &domain.SubscriptionResponse{
			IsActive:     false,
			Subscription: nil,
		}, nil
	}

	role := domain.SubscriptionRoleMember

	if subscription.BackedBy == userId {
		role = domain.SubscriptionRolePayer
	}

	members, err := s.subscriptionRepository.GetSubscriptionMembers(ctx, subscription.ID)
	if err != nil {
		return nil, err
	}

	membersResponse := make([]domain.SubscriptionMember, 0, len(members))
	for _, member := range members {
		user, err := s.userService.GetAccount(ctx, member.UserID)
		if err != nil {
			log.L(ctx).Errorf("Failed to get account details for subscription member %s: %v", member.UserID, err)
			continue
		}

		role := domain.SubscriptionRoleMember

		if subscription.BackedBy == userId {
			role = domain.SubscriptionRolePayer
		}

		if user == nil {
			membersResponse = append(membersResponse, domain.SubscriptionMember{
				UserID:          member.UserID,
				Email:           "",
				Name:            "",
				Role:            role,
				IsAccountActive: false,
			})

			continue
		}

		membersResponse = append(membersResponse, domain.SubscriptionMember{
			UserID:          member.UserID,
			Email:           user.Email,
			Name:            user.Name,
			Role:            role,
			IsAccountActive: true,
		})
	}

	return &domain.SubscriptionResponse{
		IsActive: true,
		Subscription: &domain.Subscription{
			Role:    role,
			Level:   subscription.Level,
			Till:    subscription.Till,
			Since:   subscription.Since,
			Members: membersResponse,
		},
	}, nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context) error {
	userId, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	subscription, err := s.subscriptionRepository.GetSubscriptionByMemberUserID(ctx, userId)
	if err != nil {
		return err
	}

	if subscription == nil {
		return appErr.Errorf(appErr.SubscriptionIsRequiredError, "Subscription is required")
	}

	if subscription.BackedBy != userId {
		return appErr.Errorf(appErr.PayerPermissionRequired, "You are not a payer for subscription %s", subscription.ID)
	}

	// Subscription will be removed from DB after webhook will be processed.
	return s.billingAPI.CancelSubscription(ctx, subscription.ID)
}

func (s *subscriptionService) RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID) error {
	userId, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	subscription, err := s.subscriptionRepository.GetSubscriptionByPayerUserID(ctx, userId)
	if err != nil {
		return err
	}

	if subscription == nil {
		return appErr.Errorf(appErr.SubscriptionIsRequiredError, "Subscription is required")
	}

	if subscription.BackedBy != userId {
		return appErr.Errorf(appErr.PayerPermissionRequired, "You are not a payer for subscription %s", subscription.ID)
	}

	return s.subscriptionRepository.RemoveSubscriptionMemberBySubscriptionAndUserID(ctx, subscription.ID, memberUserID)
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

	subscription, err := s.subscriptionRepository.GetSubscriptionByPayerUserID(ctx, account.ID)
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

	return s.subscriptionRepository.PutSubscription(ctx, model.Subscription{
		ID:       request.SubscriptionID,
		BackedBy: account.ID,
		Level:    subscriptionLevel,
		Till:     request.Till,
		Since:    request.Since,
	})
}

func (s *subscriptionService) RemoveSubscription(ctx context.Context, subscriptionID string) error {
	return s.subscriptionRepository.RemoveSubscription(ctx, subscriptionID)
}
