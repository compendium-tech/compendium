package service

import "github.com/seacite-tech/compendium/user-service/internal/repository"

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}
