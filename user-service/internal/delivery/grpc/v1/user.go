package grpcv1

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	myerror "github.com/compendium-tech/compendium/user-service/internal/error"
	pb "github.com/compendium-tech/compendium/user-service/internal/proto/v1"
	"github.com/compendium-tech/compendium/user-service/internal/service"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserServiceServer(userService service.UserService) UserServiceServer {
	return UserServiceServer{
		userService: userService,
	}
}

func (s UserServiceServer) Register(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, s)
	reflection.Register(server)
}

func (s UserServiceServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (_ *pb.Account, e error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				var myerr myerror.MyError
				if errors.As(err, &myerr) && myerr.ErrorType() == myerror.UserNotFoundError {
					return
				}

				e = status.Errorf(codes.Internal, "failed to get user: %v", err)
			}
		}
	}()

	if req == nil || req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	user := s.userService.GetAccount(ctx, userID)
	return &pb.Account{
		Id:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
