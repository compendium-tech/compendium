package v1

import (
	"context"

	appErr "github.com/compendium-tech/compendium/user-service/internal/error"
	pb "github.com/compendium-tech/compendium/user-service/internal/proto/v1"
	"github.com/compendium-tech/compendium/user-service/internal/service"
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

func (s *UserServiceServer) FindAccountByEmail(ctx context.Context, req *pb.FindAccountByEmailRequest) (*pb.FindAccountByEmailResponse, error) {
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email cannot be empty")
	}

	user, err := s.userService.FindAccountByEmail(ctx, req.Email)
	if err != nil {
		if err, ok := err.(appErr.AppError); ok && err.Kind() == appErr.UserNotFoundError {
			return nil, status.Errorf(codes.NotFound, "user with email %s not found", req.Email)
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.FindAccountByEmailResponse{Account: &pb.Account{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}}, nil
}
