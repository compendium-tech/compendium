package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/seacite-tech/compendium/common/pkg/log"
	"github.com/seacite-tech/compendium/common/pkg/middleware"
	"github.com/seacite-tech/compendium/user-service/internal/config"
	v1 "github.com/seacite-tech/compendium/user-service/internal/controller/v1"
	"github.com/seacite-tech/compendium/user-service/internal/email"
	"github.com/seacite-tech/compendium/user-service/internal/hash"
	"github.com/seacite-tech/compendium/user-service/internal/repository"
	"github.com/seacite-tech/compendium/user-service/internal/service"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	Config         *config.AppConfig
	PgDb           *sql.DB
	RedisClient    *redis.Client
	TokenManager   auth.TokenManager
	EmailSender    email.EmailSender
	PasswordHasher hash.PasswordHasher
}

func NewApp(deps Dependencies) *gin.Engine {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	authEmailLockRepository := repository.NewRedisAuthEmailLockRepository(deps.RedisClient)
	deviceRepository := repository.NewPgDeviceRepository(deps.PgDb)
	userRepository := repository.NewPgUserRepository(deps.PgDb)
	mfaRepository := repository.NewRedisMfaRepository(deps.RedisClient)
	refreshTokenRepository := repository.NewRedisRefreshTokenRepository(deps.RedisClient)

	authService := service.NewAuthService(
		authEmailLockRepository, deviceRepository,
		userRepository, mfaRepository, refreshTokenRepository,
		deps.EmailSender, deps.TokenManager, deps.PasswordHasher)

	r := gin.Default()
	r.Use(middleware.RequestIdMiddleware{AllowToSet: false}.Handle())
	r.Use(middleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle())
	r.Use(middleware.DefaultCors().Handle())

	v1.NewAuthController(authService).MakeRoutes(r)

	return r
}
