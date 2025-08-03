package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	errorutils "github.com/compendium-tech/compendium/common/pkg/error"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/random"

	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
)

type SubscriptionService interface {
	GetSubscription(ctx context.Context) (*domain.SubscriptionResponse, error)
	GetSubscriptionInvitationCode(ctx context.Context) (*domain.InvitationCodeResponse, error)

	UpdateSubscriptionInvitationCode(ctx context.Context) (*domain.InvitationCodeResponse, error)
	RemoveSubscriptionInvitationCode(ctx context.Context) (*domain.InvitationCodeResponse, error)

	CancelSubscription(ctx context.Context) error

	JoinCollectiveSubscription(ctx context.Context, invitationCode string) (*domain.SubscriptionResponse, error)
	RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID) error

	HandleUpdatedSubscription(ctx context.Context, request domain.HandleUpdatedSubscriptionRequest) error
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
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Getting subscription for authenticated user")

	subscription, err := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Subscription details fetched successfully")

	return s.subscriptionToResponse(ctx, userID, subscription)
}

func (s *subscriptionService) GetSubscriptionInvitationCode(ctx context.Context) (*domain.InvitationCodeResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Getting subscription invitation code for authenticated user")

	subscription, err := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, invitation code cannot be retrieved")
		return nil, myerror.New(myerror.SubscriptionIsRequiredError)
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, invitation code cannot be retrieved", subscription.ID)
		return nil, myerror.New(myerror.PayerPermissionRequired)
	}

	log.L(ctx).Info("Subscription invitation code fetched successfully")

	return &domain.InvitationCodeResponse{
		InvitationCode: subscription.InvitationCode,
	}, nil
}

func (s *subscriptionService) UpdateSubscriptionInvitationCode(ctx context.Context) (*domain.InvitationCodeResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Updating subscription invitation code for authenticated user")

	subscription, err := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, invitation code cannot be updated")
		return nil, myerror.New(myerror.SubscriptionIsRequiredError)
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, invitation code cannot be updated", subscription.ID)
		return nil, myerror.New(myerror.PayerPermissionRequired)
	}

	invitationCode, err := random.NewRandomString(12)
	if err != nil {
		return nil, err
	}

	err = s.subscriptionRepository.UpsertSubscription(ctx, model.Subscription{
		ID:             subscription.ID,
		BackedBy:       subscription.BackedBy,
		Tier:           subscription.Tier,
		InvitationCode: &invitationCode,
		Till:           subscription.Till,
		Since:          subscription.Since,
	})
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Subscription invitation code updated successfully")

	return &domain.InvitationCodeResponse{
		InvitationCode: &invitationCode,
	}, nil
}

func (s *subscriptionService) RemoveSubscriptionInvitationCode(ctx context.Context) (*domain.InvitationCodeResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Removing subscription invitation code for authenticated user")

	subscription, err := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, invitation code cannot be removed")
		return nil, myerror.New(myerror.SubscriptionIsRequiredError)
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, invitation code cannot be removed", subscription.ID)
		return nil, myerror.New(myerror.PayerPermissionRequired)
	}

	err = s.subscriptionRepository.UpsertSubscription(ctx, model.Subscription{
		ID:             subscription.ID,
		BackedBy:       subscription.BackedBy,
		Tier:           subscription.Tier,
		InvitationCode: nil,
		Till:           subscription.Till,
		Since:          subscription.Since,
	})
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Subscription invitation code removed successfully")

	return &domain.InvitationCodeResponse{
		InvitationCode: nil,
	}, nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	log.L(ctx).Info("Cancelling subscription for authenticated user")

	subscription, err := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if err != nil {
		return err
	}

	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, cannot cancel")
		return myerror.New(myerror.SubscriptionIsRequiredError)
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, cannot cancel", subscription.ID)
		return myerror.New(myerror.PayerPermissionRequired)
	}

	err = s.billingAPI.CancelSubscription(ctx, subscription.ID)
	if err != nil {
		return err
	}

	err = s.subscriptionRepository.RemoveSubscription(ctx, subscription.ID)
	if err != nil {
		return err
	}

	log.L(ctx).Info("Subscription cancellation initiated successfully")

	return nil
}

func (s *subscriptionService) JoinCollectiveSubscription(ctx context.Context, invitationCode string) (*domain.SubscriptionResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("invitationCode", invitationCode)
	logger.Info("Attempting to join collective subscription")

	subscription, err := s.subscriptionRepository.FindSubscriptionByInvitationCode(ctx, invitationCode)
	if err != nil {
		return nil, err
	}

	if subscription == nil || subscription.Tier == model.TierStudent {
		logger.Warn("Invalid invitation code or student tier subscription used to join collective subscription")
		return nil, myerror.New(myerror.InvalidSubscriptionInvitationCode)
	}

	response, err := s.subscriptionToResponse(ctx, userID, subscription)
	if err != nil {
		return nil, err
	}

	err = s.subscriptionRepository.CreateSubscriptionMember(ctx, model.SubscriptionMember{
		SubscriptionID: subscription.ID,
		UserID:         userID,
		Since:          time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully joined collective subscription")

	return response, nil
}

func (s *subscriptionService) RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	logger := log.L(ctx).WithField("memberUserId", memberUserID.String())
	logger.Info("Attempting to remove subscription member")

	subscription, err := s.subscriptionRepository.FindSubscriptionByPayerUserID(ctx, userID)
	if err != nil {
		return err
	}

	if subscription == nil {
		logger.Warn("Subscription not found for the authenticated payer, cannot remove member")
		return myerror.New(myerror.SubscriptionIsRequiredError)
	}

	if subscription.BackedBy != userID {
		logger.Warnf("Authenticated user is not the payer for subscription %s, cannot remove member", subscription.ID)
		return myerror.New(myerror.PayerPermissionRequired)
	}

	err = s.subscriptionRepository.RemoveSubscriptionMemberBySubscriptionAndUserID(ctx, subscription.ID, memberUserID)
	if err != nil {
		return err
	}

	logger.Info("Subscription member removed successfully")

	return nil
}

func (s *subscriptionService) HandleUpdatedSubscription(ctx context.Context, request domain.HandleUpdatedSubscriptionRequest) (finalErr error) {
	logger := log.L(ctx).WithField("subscriptionId", request.SubscriptionID).WithField("userId", request.UserID)
	logger.Info("Upserting subscription")

	if len(request.Items) == 0 {
		logger.Warn("No items in upsert subscription request")
		return myerror.NewWithReason(myerror.RequestValidationError, "No items in request")
	}

	lock, err := s.billingLockRepository.ObtainLock(ctx, request.UserID)
	if err != nil {
		logger.Errorf("Failed to obtain billing lock for customer %s: %v", request.UserID, err)
		return err
	}

	defer errorutils.DeferTryWithContext(ctx, &finalErr, lock.Release)

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
				isCanceled, err := s.billingAPI.IsSubscriptionCanceled(ctx, request.SubscriptionID)
				if err != nil {
					return err
				}

				if !isCanceled {
					err = s.billingAPI.CancelSubscription(ctx, request.SubscriptionID)
					if err != nil {
						return err
					}
				}

				return myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("Unknown product ID: %s", item.ProductID))
			}
			continue
		}

		// To maintain a single active subscription per payer within this service,
		// we first cancel any existing subscription before upserting a new one.
		//
		// If the cancellation fails, the external billing service will retry the request,
		// ensuring the database isn't updated with a new subscription until the previous one is successfully canceled.
		existingSubscription, err := s.subscriptionRepository.FindSubscriptionByPayerUserID(ctx, request.UserID)
		if err != nil {
			return err
		}

		if existingSubscription != nil {
			itemLogger.Infof("Existing subscription %s found, cancelling it before upserting new item", existingSubscription.ID)

			isCanceled, err := s.billingAPI.IsSubscriptionCanceled(ctx, existingSubscription.ID)
			if err != nil {
				return err
			}

			if !isCanceled {
				err = s.billingAPI.CancelSubscription(ctx, existingSubscription.ID)
				if err != nil {
					return err
				}
			}
		}

		err = s.subscriptionRepository.UpsertSubscription(ctx, model.Subscription{
			ID:       request.SubscriptionID,
			BackedBy: request.UserID,
			Tier:     subscriptionLevel,
			Till:     request.Till,
			Since:    request.Since,
		})
		if err != nil {
			return err
		}

		itemLogger.Info("Subscription item processed successfully")
	}

	logger.Info("All subscription items processed, upsert completed successfully")

	return nil
}

func (s *subscriptionService) RemoveSubscription(ctx context.Context, subscriptionID string) error {
	logger := log.L(ctx).WithField("subscriptionId", subscriptionID)
	logger.Info("Removing subscription")

	err := s.subscriptionRepository.RemoveSubscription(ctx, subscriptionID)
	if err != nil {
		return err
	}

	logger.Info("Subscription removed successfully")

	return nil
}

func (s *subscriptionService) subscriptionToResponse(ctx context.Context, userID uuid.UUID, subscription *model.Subscription) (*domain.SubscriptionResponse, error) {
	if subscription == nil {
		return &domain.SubscriptionResponse{
			IsActive:     false,
			Subscription: nil,
		}, nil
	}

	role := domain.SubscriptionRoleMember

	if subscription.BackedBy == userID {
		role = domain.SubscriptionRolePayer
	}

	logger := log.L(ctx).WithField("subscriptionId", subscription.ID)
	logger.Debug("Converting subscription to response format")

	members, err := s.subscriptionRepository.GetSubscriptionMembers(ctx, subscription.ID)
	if err != nil {
		logger.Errorf("Failed to get subscription members: %v", err)
		return nil, err
	}

	membersResponse := make([]domain.SubscriptionMember, 0, len(members))
	for _, member := range members {
		user, err := s.userService.GetAccount(ctx, member.UserID)
		if err != nil {
			logger.Errorf("Failed to get account details for subscription member %s: %v", member.UserID, err)
			continue
		}

		memberRole := domain.SubscriptionRoleMember
		if subscription.BackedBy == member.UserID {
			memberRole = domain.SubscriptionRolePayer
		}

		if user == nil {
			membersResponse = append(membersResponse, domain.SubscriptionMember{
				UserID:          member.UserID,
				Email:           "",
				Name:            "",
				Role:            memberRole,
				IsAccountActive: false,
			})
			logger.Debugf("Added inactive member %s to response", member.UserID)
			continue
		}

		membersResponse = append(membersResponse, domain.SubscriptionMember{
			UserID:          member.UserID,
			Email:           user.Email,
			Name:            user.Name,
			Role:            memberRole,
			IsAccountActive: true,
		})
		logger.Debugf("Added active member %s to response", member.UserID)
	}

	logger.Debug("Successfully converted subscription to response format")

	return &domain.SubscriptionResponse{
		IsActive: true,
		Subscription: &domain.Subscription{
			Role:    role,
			Tier:    subscription.Tier,
			Till:    subscription.Till,
			Since:   subscription.Since,
			Members: membersResponse,
		},
	}, nil
}
