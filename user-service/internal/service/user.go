package service

import "github.com/seacite-tech/compendium/user-service/internal/repository"

type UserService interface{}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userServiceImpl {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}
