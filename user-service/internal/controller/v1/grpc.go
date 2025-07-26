package v1

import (
	"context"

	appErr "github.com/compendium-tech/compendium/user-service/internal/error"
	pb "github.com/compendium-tech/compendium/user-service/internal/proto/v1"
	"github.com/compendium-tech/compendium/user-service/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserServiceServer(userService service.UserService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

func (s *UserServiceServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	user, err := s.userService.GetAccount(ctx, userID)
	if err != nil {
		if err, ok := err.(appErr.AppError); ok && err.Kind() == appErr.UserNotFoundError {
			return nil, nil
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.Account{
		Id:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
