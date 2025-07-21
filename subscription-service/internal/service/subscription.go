package service

import (
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
)

type SubscriptionService interface {
	PutSubscription(request domain.PutSubscriptionRequest) error
}

type subscriptionService struct {
	subscriptionRepository repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepository repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *subscriptionService) PutSubscription(request domain.PutSubscriptionRequest) error {
	return s.subscriptionRepository.PutSubscription(model.Subscription{
		UserID:            request.UserID,
		SubscriptionLevel: request.SubscriptionLevel,
		Till:              request.Till,
		Since:             request.Since,
	})
}
