package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
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
	GetSubscription(ctx context.Context) domain.SubscriptionResponse
	GetSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse

	UpdateSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse
	RemoveSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse

	CancelSubscription(ctx context.Context)

	JoinCollectiveSubscription(ctx context.Context, invitationCode string) domain.SubscriptionResponse
	RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID)

	HandleUpdatedSubscription(ctx context.Context, request domain.HandleUpdatedSubscriptionRequest)
	RemoveSubscription(ctx context.Context, subscriptionID string)
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

func (s *subscriptionService) GetSubscription(ctx context.Context) domain.SubscriptionResponse {
	userID := auth.GetUserID(ctx)
	log.L(ctx).Info("Getting subscription for authenticated user")

	subscription := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	log.L(ctx).Info("Subscription details fetched successfully")

	return s.subscriptionToResponse(ctx, userID, subscription)
}

func (s *subscriptionService) GetSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse {
	userID := auth.GetUserID(ctx)
	log.L(ctx).Info("Getting subscription invitation code for authenticated user")

	subscription := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, invitation code cannot be retrieved")
		myerror.New(myerror.SubscriptionIsRequiredError).Throw()
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, invitation code cannot be retrieved", subscription.ID)
		myerror.New(myerror.PayerPermissionRequired).Throw()
	}

	log.L(ctx).Info("Subscription invitation code fetched successfully")

	return domain.InvitationCodeResponse{
		InvitationCode: subscription.InvitationCode,
	}
}

func (s *subscriptionService) UpdateSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse {
	userID := auth.GetUserID(ctx)
	log.L(ctx).Info("Updating subscription invitation code for authenticated user")

	subscription := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, invitation code cannot be updated")
		myerror.New(myerror.SubscriptionIsRequiredError).Throw()
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, invitation code cannot be updated", subscription.ID)
		myerror.New(myerror.PayerPermissionRequired).Throw()
	}

	invitationCode := random.NewRandomString(12)
	s.subscriptionRepository.UpsertSubscription(ctx, model.Subscription{
		ID:             subscription.ID,
		BackedBy:       subscription.BackedBy,
		Tier:           subscription.Tier,
		InvitationCode: &invitationCode,
		Till:           subscription.Till,
		Since:          subscription.Since,
	})

	log.L(ctx).Info("Subscription invitation code updated successfully")

	return domain.InvitationCodeResponse{
		InvitationCode: &invitationCode,
	}
}

func (s *subscriptionService) RemoveSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse {
	userID := auth.GetUserID(ctx)
	log.L(ctx).Info("Removing subscription invitation code for authenticated user")

	subscription := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, invitation code cannot be removed")
		myerror.New(myerror.SubscriptionIsRequiredError).Throw()
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, invitation code cannot be removed", subscription.ID)
		myerror.New(myerror.PayerPermissionRequired).Throw()
	}

	s.subscriptionRepository.UpsertSubscription(ctx, model.Subscription{
		ID:             subscription.ID,
		BackedBy:       subscription.BackedBy,
		Tier:           subscription.Tier,
		InvitationCode: nil,
		Till:           subscription.Till,
		Since:          subscription.Since,
	})

	log.L(ctx).Info("Subscription invitation code removed successfully")

	return domain.InvitationCodeResponse{
		InvitationCode: nil,
	}
}

func (s *subscriptionService) CancelSubscription(ctx context.Context) {
	userID := auth.GetUserID(ctx)
	log.L(ctx).Info("Cancelling subscription for authenticated user")

	subscription := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	if subscription == nil {
		log.L(ctx).Warn("Subscription not found for the authenticated user, cannot cancel")
		myerror.New(myerror.SubscriptionIsRequiredError).Throw()
	}

	if subscription.BackedBy != userID {
		log.L(ctx).Warnf("Authenticated user is not the payer for subscription %s, cannot cancel", subscription.ID)
		myerror.New(myerror.PayerPermissionRequired).Throw()
	}

	s.billingAPI.CancelSubscription(ctx, subscription.ID)
	s.subscriptionRepository.RemoveSubscription(ctx, subscription.ID)

	log.L(ctx).Info("Subscription cancellation initiated successfully")
}

func (s *subscriptionService) JoinCollectiveSubscription(ctx context.Context, invitationCode string) domain.SubscriptionResponse {
	userID := auth.GetUserID(ctx)
	logger := log.L(ctx).WithField("invitationCode", invitationCode)
	logger.Info("Attempting to join collective subscription")

	subscription := s.subscriptionRepository.FindSubscriptionByInvitationCode(ctx, invitationCode)
	if subscription == nil || subscription.Tier == model.TierStudent {
		logger.Warn("Invalid invitation code or student tier subscription used to join collective subscription")
		myerror.New(myerror.InvalidSubscriptionInvitationCode).Throw()
	}

	response := s.subscriptionToResponse(ctx, userID, subscription)
	s.subscriptionRepository.CreateSubscriptionMemberAndCheckMemberCount(ctx, model.SubscriptionMember{
		SubscriptionID: subscription.ID,
		UserID:         userID,
		Since:          time.Now().UTC(),
	}, func(memberCount uint) error {
		var tierMemberLimits = map[model.Tier]uint{
			model.TierStudent:   1,
			model.TierTeam:      10,
			model.TierCommunity: 40,
		}

		limit := tierMemberLimits[subscription.Tier]

		if memberCount > limit {
			return myerror.New(myerror.InvalidSubscriptionInvitationCode)
		}

		return nil
	})

	logger.Info("Successfully joined collective subscription")

	return response
}

func (s *subscriptionService) RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID) {
	userID := auth.GetUserID(ctx)
	logger := log.L(ctx).WithField("memberUserId", memberUserID.String())
	logger.Info("Attempting to remove subscription member")

	subscription := s.subscriptionRepository.FindSubscriptionByPayerUserID(ctx, userID)
	if subscription == nil {
		logger.Warn("Subscription not found for the authenticated payer, cannot remove member")
		myerror.New(myerror.SubscriptionIsRequiredError).Throw()
	}

	if subscription.BackedBy != userID {
		logger.Warnf("Authenticated user is not the payer for subscription %s, cannot remove member", subscription.ID)
		myerror.New(myerror.PayerPermissionRequired).Throw()
	}

	s.subscriptionRepository.RemoveSubscriptionMemberBySubscriptionAndUserID(ctx, subscription.ID, memberUserID)
	logger.Info("Subscription member removed successfully")
}

func (s *subscriptionService) HandleUpdatedSubscription(ctx context.Context, request domain.HandleUpdatedSubscriptionRequest) {
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

func (s *subscriptionService) RemoveSubscription(ctx context.Context, subscriptionID string) {
	logger := log.L(ctx).WithField("subscriptionId", subscriptionID)
	logger.Info("Removing subscription")

	s.subscriptionRepository.RemoveSubscription(ctx, subscriptionID)
	logger.Info("Subscription removed successfully")
}

func (s *subscriptionService) subscriptionToResponse(ctx context.Context, userID uuid.UUID, subscription *model.Subscription) domain.SubscriptionResponse {
	if subscription == nil {
		return domain.SubscriptionResponse{
			IsActive:     false,
			Subscription: nil,
		}
	}

	role := domain.SubscriptionRoleMember

	if subscription.BackedBy == userID {
		role = domain.SubscriptionRolePayer
	}

	logger := log.L(ctx).WithField("subscriptionId", subscription.ID)
	logger.Debug("Converting subscription to response format")

	members := s.subscriptionRepository.GetSubscriptionMembers(ctx, subscription.ID)

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

	return domain.SubscriptionResponse{
		IsActive: true,
		Subscription: &domain.Subscription{
			Role:    role,
			Tier:    subscription.Tier,
			Till:    subscription.Till,
			Since:   subscription.Since,
			Members: membersResponse,
		},
	}
}
