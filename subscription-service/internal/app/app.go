package app

import (
	"database/sql"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/adslmgrv/compendium/subscription-service/internal/config"
	"github.com/adslmgrv/compendium/subscription-service/internal/controller"
	"github.com/adslmgrv/compendium/subscription-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/seacite-tech/compendium/common/pkg/log"
	commonMiddleware "github.com/seacite-tech/compendium/common/pkg/middleware"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	Config       *config.AppConfig
	PgDb         *sql.DB
	RedisClient  *redis.Client
	TokenManager auth.TokenManager
}

func NewApp(deps Dependencies) *gin.Engine {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	subscriptionRepository := repository.NewPgSubscriptionRepository(deps.PgDb)

	r := gin.Default()
	r.Use(commonMiddleware.RequestIdMiddleware{AllowToSet: false}.Handle)
	r.Use(auth.AuthMiddleware{TokenManager: deps.TokenManager}.Handle)
	r.Use(commonMiddleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle)
	r.Use(commonMiddleware.DefaultCors().Handle)

	controller.NewPaddleWebhookController(
		subscriptionRepository,
		*paddle.NewWebhookVerifier(deps.Config.PaddleWebhookSecret)).MakeRoutes(r)

	return r
}
