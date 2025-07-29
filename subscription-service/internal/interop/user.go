package interop

import (
	"context"
	"fmt"

	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	pb "github.com/compendium-tech/compendium/subscription-service/internal/proto/v1"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userServiceGrpcClient struct {
	client pb.UserServiceClient
}

type UserService interface {
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error)
}

func NewGrpcUserServiceClient(target string) (UserService, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	c := pb.NewUserServiceClient(conn)

	return &userServiceGrpcClient{
		client: c,
	}, nil
}

func (u *userServiceGrpcClient) GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	req := &pb.GetAccountRequest{
		Id: id.String(),
	}

	resp, err := u.client.GetAccount(ctx, req)
	if err != nil {
		return nil, tracerr.Errorf("failed to get account: %w", err)
	}

	accountID, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, tracerr.Errorf("invalid account ID format: %w", err)
	}

	account := &domain.Account{
		ID:    accountID,
		Name:  resp.Name,
		Email: resp.Email,
	}

	return account, nil
}
