package app

import (
	"database/sql"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/middleware"

	"github.com/compendium-tech/compendium/subscription-service/internal/billing"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	httpv1 "github.com/compendium-tech/compendium/subscription-service/internal/delivery/http/v1"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/repository"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
)

type Dependencies struct {
	Config          *config.AppConfig
	PgDB            *sql.DB
	RedisClient     *redis.Client
	TokenManager    auth.TokenManager
	UserService     interop.UserService
	PaddleAPIClient paddle.SDK
}

func NewApp(deps Dependencies) *gin.Engine {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "subscription-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	billingAPI := billing.NewPaddleBillingAPI(deps.PaddleAPIClient)
	subscriptionRepository := repository.NewPgSubscriptionRepository(deps.PgDB)
	billingLockRepository := repository.NewRedisBillingLockRepository(deps.RedisClient)
	subscriptionService := service.NewSubscriptionService(
		billingAPI, deps.Config.ProductIDs, deps.UserService,
		billingLockRepository, subscriptionRepository)

	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware{AllowToSet: false}.Handle)
	r.Use(auth.Middleware{TokenManager: deps.TokenManager}.Handle)
	r.Use(middleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle)

	httpv1.NewBillingWebhookController(
		subscriptionService,
		paddle.NewWebhookVerifier(deps.Config.PaddleWebhookSecret)).MakeRoutes(r)
	httpv1.NewSubscriptionController(subscriptionService).MakeRoutes(r)

	return r
}
