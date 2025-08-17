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
	GetAccount(ctx context.Context, id uuid.UUID) domain.AccountResponse
	FindAccountByEmail(ctx context.Context, email string) domain.AccountResponse

	GetAccountAsAuthenticatedUser(ctx context.Context) domain.AccountResponse
	UpdateAccountAsAuthenticatedUser(ctx context.Context, request domain.UpdateAccount) domain.AccountResponse
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) GetAccount(ctx context.Context, id uuid.UUID) domain.AccountResponse {
	logger := log.L(ctx).WithField("userId", id.String())
	logger.Info("Getting user account details by ID")

	user := u.userRepository.GetUser(ctx, id)

	if user == nil {
		logger.Warn("User account not found for the provided ID")
		myerror.New(myerror.UserNotFoundError).Throw()
	}

	logger.Info("User account details fetched successfully")

	return domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (u *userService) FindAccountByEmail(ctx context.Context, email string) domain.AccountResponse {
	logger := log.L(ctx).WithField("email", email)
	logger.Info("Finding user account by email")

	user := u.userRepository.FindUserByEmail(ctx, email)

	if user == nil {
		logger.Warn("User account not found for the provided email")
		myerror.New(myerror.UserNotFoundError).Throw()
	}

	logger.Info("User account found successfully")

	return domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (u *userService) GetAccountAsAuthenticatedUser(ctx context.Context) domain.AccountResponse {
	userID := auth.GetUserID(ctx)

	log.L(ctx).Info("Fetching authenticated user account details")

	user := u.userRepository.GetUser(ctx, userID)
	if user == nil {
		log.L(ctx).Warn("Information about authenticated user not found, perhaps session is invalid?")
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	log.L(ctx).Info("Account details fetched successfully")

	return domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (u *userService) UpdateAccountAsAuthenticatedUser(ctx context.Context, request domain.UpdateAccount) domain.AccountResponse {
	userID := auth.GetUserID(ctx)

	logger := log.L(ctx).WithField("newName", request.Name)
	logger.Info("Updating authenticated user account details")

	user := u.userRepository.UpdateUserName(ctx, userID, request.Name)

	logger.Info("Account details updated successfully")

	return domain.AccountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
