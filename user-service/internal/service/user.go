package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"

	"github.com/compendium-tech/compendium/user-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/user-service/internal/error"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
)

// UserService defines the interface for managing user accounts.
//
// Methods suffixed with AsAuthenticatedUser are designed for use within
// authenticated API endpoints. For these, the user's identity is automatically
// derived from the context, populated by an authentication middleware.
//
// Methods without this suffix (like FindAccountByEmail) are
// intended for interoperability with other microservices.
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
	logger := log.L(ctx).WithField("userId", id.String())
	logger.Info("Getting user account details by ID")

	user, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		logger.Warn("User account not found for the provided ID")
		return nil, myerror.New(myerror.UserNotFoundError)
	}

	logger.Info("User account details fetched successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) FindAccountByEmail(ctx context.Context, email string) (*domain.AccountResponse, error) {
	logger := log.L(ctx).WithField("email", email)
	logger.Info("Finding user account by email")

	user, err := u.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		logger.Warn("User account not found for the provided email")
		return nil, myerror.New(myerror.UserNotFoundError)
	}

	logger.Info("User account found successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) GetAccountAsAuthenticatedUser(ctx context.Context) (*domain.AccountResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("Fetching authenticated user account details")

	user, err := u.userRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Warn("Information about authenticated user not found, perhaps session is invalid?")
		return nil, myerror.New(myerror.InvalidSessionError)
	}

	log.L(ctx).Info("Account details fetched successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userService) UpdateAccountAsAuthenticatedUser(ctx context.Context, request domain.UpdateAccount) (*domain.AccountResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("newName", request.Name)
	logger.Info("Updating authenticated user account details")

	user, err := u.userRepository.UpdateUserName(ctx, userID, request.Name)
	if err != nil {
		return nil, err
	}

	logger.Info("Account details updated successfully")

	return &domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
