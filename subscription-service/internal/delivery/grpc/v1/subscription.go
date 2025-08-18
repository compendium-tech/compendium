package grpcv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	pb "github.com/compendium-tech/compendium/subscription-service/internal/proto/v1"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/google/uuid"
)

type SubscriptionServiceServer struct {
	pb.UnimplementedSubscriptionServiceServer
	subscriptionService service.SubscriptionService
}

func NewSubscriptionServiceServer(subscriptionService service.SubscriptionService) SubscriptionServiceServer {
	return SubscriptionServiceServer{
		subscriptionService: subscriptionService,
	}
}

func (s SubscriptionServiceServer) Register(server *grpc.Server) {
	pb.RegisterSubscriptionServiceServer(server, s)
	reflection.Register(server)
}

func (s SubscriptionServiceServer) GetSubscriptionTier(ctx context.Context, req *pb.GetSubscriptionTierRequest) (_ *pb.GetSubscriptionTierResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				e = status.Errorf(codes.Internal, "failed to get subscription tier: %v", err)
			}
		}
	}()

	if req == nil || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	tier := s.subscriptionService.GetSubscriptionTierByMemberUserID(ctx, userID)
	if tier == nil {
		return &pb.GetSubscriptionTierResponse{Tier: pb.SubscriptionTier_NONE}, nil
	}

	var pbTier pb.SubscriptionTier
	switch *tier {
	case model.TierStudent:
		pbTier = pb.SubscriptionTier_STUDENT
	case model.TierTeam:
		pbTier = pb.SubscriptionTier_TEAM
	case model.TierCommunity:
		pbTier = pb.SubscriptionTier_COMMUNITY
	}

	return &pb.GetSubscriptionTierResponse{
		Tier: pbTier,
	}, nil
}
