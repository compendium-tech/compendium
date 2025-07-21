package app

import (
	"database/sql"

	"github.com/compendium-tech/compendium/common/pkg/log"
	commonMiddleware "github.com/compendium-tech/compendium/common/pkg/middleware"
	emailDelivery "github.com/compendium-tech/compendium/email-delivery-service/pkg/email"
	"github.com/compendium-tech/compendium/user-service/internal/config"
	v1 "github.com/compendium-tech/compendium/user-service/internal/controller/v1"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/compendium-tech/compendium/user-service/internal/service"
	"github.com/compendium-tech/compendium/user-service/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	Config              config.AppConfig
	PgDb                *sql.DB
	RedisClient         *redis.Client
	TokenManager        auth.TokenManager
	EmailSender         emailDelivery.EmailSender
	EmailMessageBuilder email.EmailMessageBuilder
	PasswordHasher      hash.PasswordHasher
}

func NewApp(deps Dependencies) *gin.Engine {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	authEmailLockRepository := repository.NewRedisAuthLockRepository(deps.RedisClient)
	deviceRepository := repository.NewPgDeviceRepository(deps.PgDb)
	userRepository := repository.NewPgUserRepository(deps.PgDb)
	mfaRepository := repository.NewRedisMfaRepository(deps.RedisClient)
	refreshTokenRepository := repository.NewRedisRefreshTokenRepository(deps.RedisClient)

	authService := service.NewAuthService(
		authEmailLockRepository, deviceRepository,
		userRepository, mfaRepository, refreshTokenRepository,
		deps.EmailSender, deps.EmailMessageBuilder, deps.TokenManager, deps.PasswordHasher)
	userService := service.NewUserService(userRepository)

	r := gin.Default()
	r.Use(commonMiddleware.RequestIdMiddleware{AllowToSet: false}.Handle)
	r.Use(auth.AuthMiddleware{TokenManager: deps.TokenManager}.Handle)
	r.Use(commonMiddleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle)
	r.Use(commonMiddleware.DefaultCors().Handle)

	v1.NewAuthController(authService).MakeRoutes(r)
	v1.NewUserController(userService).MakeRoutes(r)

	return r
}
