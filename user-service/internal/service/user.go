package service

import (
	"context"

	"github.com/seacite-tech/compendium/common/pkg/auth"
	log "github.com/seacite-tech/compendium/common/pkg/log"
	"github.com/seacite-tech/compendium/user-service/internal/domain"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/repository"
)

type UserService interface {
	GetAccount(ctx context.Context) (*domain.AccountResponse, error)
	UpdateAccount(ctx context.Context, request domain.UpdateAccount) (*domain.AccountResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) GetAccount(ctx context.Context) (*domain.AccountResponse, error) {
	log.L(ctx).Info("Getting authenticated user account details")

	userId, err := auth.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Fetching authenticated user account details in database")
	user, err := u.userRepository.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Warn("Information about authenticated user not found, perhaphs session is invalid?")

		return nil, appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	log.L(ctx).Info("Account details fetched successfully")

	return &domain.AccountResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) UpdateAccount(ctx context.Context, request domain.UpdateAccount) (*domain.AccountResponse, error) {
	log.L(ctx).Info("Updating authenticated user account details")

	userId, err := auth.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Updating authenticated user account details in database")
	user, err := u.userRepository.UpdateName(ctx, userId, request.Name)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Account details updated successfully")

	return &domain.AccountResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
