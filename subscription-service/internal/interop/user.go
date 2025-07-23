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
	FindAccountByEmail(ctx context.Context, email string) (*domain.Account, error)
}

func NewGrpcUserServiceClient(target string) (UserService, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to gRPC server: %v", err)
		return nil, err
	}

	c := pb.NewUserServiceClient(conn)

	return &userServiceGrpcClient{
		client: c,
	}, nil
}

func (u *userServiceGrpcClient) FindAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	req := &pb.FindAccountByEmailRequest{
		Email: email,
	}

	resp, err := u.client.FindAccountByEmail(ctx, req)
	if err != nil {
		return nil, tracerr.Errorf("failed to get account: %w", err)
	}

	id, err := uuid.Parse(resp.Account.ID)
	if err != nil {
		return nil, tracerr.Errorf("invalid account ID format: %w", err)
	}

	account := &domain.Account{
		ID:    id,
		Name:  resp.Account.Name,
		Email: resp.Account.Email,
	}

	return account, nil
}
