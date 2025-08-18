package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/random"

	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
)

type SubscriptionService interface {
	GetSubscriptionTierByMemberUserID(ctx context.Context, userID uuid.UUID) *model.Tier
	GetSubscription(ctx context.Context) domain.SubscriptionResponse
	GetSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse
	UpdateSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse
	RemoveSubscriptionInvitationCode(ctx context.Context) domain.InvitationCodeResponse
	CancelSubscription(ctx context.Context)
	JoinCollectiveSubscription(ctx context.Context, invitationCode string) domain.SubscriptionResponse
	RemoveSubscriptionMember(ctx context.Context, memberUserID uuid.UUID)
}

type subscriptionService struct {
	billingAPI             billing.BillingAPI
	userService            interop.UserService
	billingLockRepository  repository.BillingLockRepository
	subscriptionRepository repository.SubscriptionRepository
}

func NewSubscriptionService(
	billingAPI billing.BillingAPI,
	userService interop.UserService,
	billingLockRepository repository.BillingLockRepository,
	subscriptionRepository repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		billingAPI:             billingAPI,
		userService:            userService,
		billingLockRepository:  billingLockRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *subscriptionService) GetSubscriptionTierByMemberUserID(ctx context.Context, userID uuid.UUID) *model.Tier {
	logger := log.L(ctx).WithField("userId", userID)
	logger.Info("Getting subscription tier")

	subscription := s.subscriptionRepository.FindSubscriptionByMemberUserID(ctx, userID)
	logger.Info("Subscription details fetched successfully")

	if subscription == nil {
		return nil
	}

	return &subscription.Tier
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
