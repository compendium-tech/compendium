package service

import (
	"context"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	log "github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/user-service/internal/domain"
	appErr "github.com/compendium-tech/compendium/user-service/internal/error"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.AccountResponse, error)
	FindAccountByEmail(ctx context.Context, email string) (*domain.AccountResponse, error)

	GetAccountAsAuthenticatedUser(ctx context.Context) (*domain.AccountResponse, error)
	UpdateAccountAsAuthenticatedUser(ctx context.Context, request domain.UpdateAccount) (*domain.AccountResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) GetAccount(ctx context.Context, id uuid.UUID) (*domain.AccountResponse, error) {
	log.L(ctx).Info("Getting user account details by ID")

	user, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Warn("User account not found for the provided ID")

		return nil, appErr.Errorf(appErr.UserNotFoundError, "User not found")
	}

	log.L(ctx).Info("User account details fetched successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) GetAccountAsAuthenticatedUser(ctx context.Context) (*domain.AccountResponse, error) {
	log.L(ctx).Info("Getting authenticated user account details")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Fetching authenticated user account details in database")
	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Warn("Information about authenticated user not found, perhaphs session is invalid?")

		return nil, appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	log.L(ctx).Info("Account details fetched successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) FindAccountByEmail(ctx context.Context, email string) (*domain.AccountResponse, error) {
	log.L(ctx).Info("Finding user account by email")

	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Warn("User account not found for the provided email")

		return nil, appErr.Errorf(appErr.UserNotFoundError, "User not found")
	}

	log.L(ctx).Info("User account found successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) UpdateAccountAsAuthenticatedUser(ctx context.Context, request domain.UpdateAccount) (*domain.AccountResponse, error) {
	log.L(ctx).Info("Updating authenticated user account details")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Updating authenticated user account details in database")
	user, err := u.userRepository.UpdateName(ctx, userID, request.Name)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Account details updated successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
