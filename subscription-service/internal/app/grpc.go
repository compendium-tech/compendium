package app

import (
	"database/sql"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/compendium-tech/compendium/common/pkg/log"
	netapp "github.com/compendium-tech/compendium/common/pkg/net"

	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	grpcv1 "github.com/compendium-tech/compendium/subscription-service/internal/delivery/grpc/v1"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
)

type GrpcAppDependencies struct {
	Config          config.GinAppConfig
	PgDB            *sql.DB
	RedisClient     *redis.Client
	UserService     interop.UserService
	PaddleAPIClient paddle.SDK
}

func NewGrpcApp(deps GrpcAppDependencies) netapp.GrpcApp {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "subscription-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	billingAPI := billing.NewPaddleBillingAPI(deps.PaddleAPIClient)
	subscriptionRepository := repository.NewPgSubscriptionRepository(deps.PgDB)
	billingLockRepository := repository.NewRedisBillingLockRepository(deps.RedisClient)
	subscriptionService := service.NewSubscriptionService(
		billingAPI, deps.UserService,
		billingLockRepository, subscriptionRepository)

	grpcServer := grpc.NewServer()
	grpcv1.NewSubscriptionServiceServer(subscriptionService).Register(grpcServer)

	return netapp.NewGrpcApp(grpcServer)
}
