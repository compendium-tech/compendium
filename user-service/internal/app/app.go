package app

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/middleware"
	"github.com/compendium-tech/compendium/common/pkg/ratelimit"
	"github.com/compendium-tech/compendium/user-service/internal/config"
	v1 "github.com/compendium-tech/compendium/user-service/internal/controller/v1"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	"github.com/compendium-tech/compendium/user-service/internal/geoip"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	pbv1 "github.com/compendium-tech/compendium/user-service/internal/proto/v1"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/compendium-tech/compendium/user-service/internal/service"
	"github.com/compendium-tech/compendium/user-service/internal/ua"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GinAppDependencies struct {
	Config              config.GinAppConfig
	PgDB                *sql.DB
	RedisClient         *redis.Client
	TokenManager        auth.TokenManager
	EmailSender         email.Sender
	EmailMessageBuilder email.MessageBuilder
	PasswordHasher      hash.PasswordHasher
	GeoIP               geoip.GeoIP
	UserAgentParser     ua.UserAgentParser
}

func NewGinApp(deps GinAppDependencies) *gin.Engine {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	rateLimiter := ratelimit.NewRedisRateLimiter(deps.RedisClient)
	authEmailLockRepository := repository.NewRedisAuthLockRepository(deps.RedisClient)
	deviceRepository := repository.NewPgTrustedDeviceRepository(deps.PgDB)
	userRepository := repository.NewPgUserRepository(deps.PgDB)
	mfaRepository := repository.NewRedisMfaRepository(deps.RedisClient)
	refreshTokenRepository := repository.NewRedisRefreshTokenRepository(deps.RedisClient)

	authService := service.NewAuthService(
		authEmailLockRepository, deviceRepository,
		userRepository, mfaRepository, refreshTokenRepository,
		deps.EmailSender, deps.EmailMessageBuilder,
		deps.GeoIP, deps.UserAgentParser,
		deps.TokenManager, deps.PasswordHasher, rateLimiter)
	userService := service.NewUserService(userRepository)

	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware{AllowToSet: false}.Handle)
	r.Use(auth.Middleware{TokenManager: deps.TokenManager}.Handle)
	r.Use(middleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle)

	v1.NewAuthController(authService).MakeRoutes(r)
	v1.NewUserController(userService).MakeRoutes(r)

	return r
}

type GrpcApp struct {
	server *grpc.Server
	port   uint16
}

func (g *GrpcApp) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logrus.Infof("Starting gRPC server on :%d", g.port)
	return g.server.Serve(lis)
}

type GrpcAppDependencies struct {
	Config config.GrpcAppConfig
	PgDB   *sql.DB
}

func NewGrpcApp(deps GrpcAppDependencies) *GrpcApp {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	userRepository := repository.NewPgUserRepository(deps.PgDB)
	userService := service.NewUserService(userRepository)

	grpcServer := grpc.NewServer()
	pbv1.RegisterUserServiceServer(grpcServer, v1.NewUserServiceServer(userService))
	reflection.Register(grpcServer)

	return &GrpcApp{server: grpcServer, port: deps.Config.GrpcPort}
}
