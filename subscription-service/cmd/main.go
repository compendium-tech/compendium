package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/joho/godotenv"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	netapp "github.com/compendium-tech/compendium/common/pkg/net"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/redis"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/compendium-tech/compendium/subscription-service/internal/app"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
)

func main() {
	appMode := flag.String("mode", "", "Specify the application mode: 'http' for Gin app, 'grpc' for gRPC app and 'webhook' for webhook Gin app")
	flag.Parse()

	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	var app netapp.App
	switch *appMode {
	case "http":
		app = createHttpApp(ctx)
	case "grpc":
		app = createGrpcApp(ctx)
	case "webhook":
		app = createWebhookApp(ctx)
	default:
		fmt.Printf("Invalid application mode specified: %s. Please use 'http' or 'grpc'.\n", *appMode)
	}

	if app == nil {
		return
	}

	err = app.Run()
	if err != nil {
		fmt.Printf("Failed to start user service, cause: %v\n", err)
	}
}

func createHttpApp(ctx context.Context) netapp.App {
	cfg := config.LoadGinAppConfig()

	tokenManager, err := auth.NewJwtBasedTokenManager(cfg.JwtSingingKey)
	if err != nil {
		fmt.Printf("Failed to initialize token manager, cause: %s", err)
		return nil
	}

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s", err)
		return nil
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		fmt.Printf("Failed to connect to Redis, cause: %s", err)
		return nil
	}

	userService, err := interop.NewGrpcUserServiceClient(cfg.GrpcUserServiceClientTarget)
	if err != nil {
		fmt.Printf("Failed to initialize user service, cause: %s", err)
		return nil
	}

	paddleAPIClient, err := paddle.New(cfg.PaddleAPIKey, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		fmt.Printf("Failed to initialize Paddle API client, cause: %s", err)
		return nil
	}

	return app.NewGinApp(app.GinAppDependencies{
		PgDB:            pgDB,
		RedisClient:     redisClient,
		Config:          cfg,
		TokenManager:    tokenManager,
		UserService:     userService,
		PaddleAPIClient: *paddleAPIClient,
	})
}

func createGrpcApp(ctx context.Context) netapp.App {
	cfg := config.LoadGinAppConfig()

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s", err)
		return nil
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		fmt.Printf("Failed to connect to Redis, cause: %s", err)
		return nil
	}

	userService, err := interop.NewGrpcUserServiceClient(cfg.GrpcUserServiceClientTarget)
	if err != nil {
		fmt.Printf("Failed to initialize user service, cause: %s", err)
		return nil
	}

	paddleAPIClient, err := paddle.New(cfg.PaddleAPIKey, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		fmt.Printf("Failed to initialize Paddle API client, cause: %s", err)
		return nil
	}

	return app.NewGrpcApp(app.GrpcAppDependencies{
		PgDB:            pgDB,
		RedisClient:     redisClient,
		Config:          cfg,
		UserService:     userService,
		PaddleAPIClient: *paddleAPIClient,
	})
}

func createWebhookApp(ctx context.Context) netapp.App {
	cfg := config.LoadGinWebhookAppConfig()

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s", err)
		return nil
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		fmt.Printf("Failed to connect to Redis, cause: %s", err)
		return nil
	}

	paddleAPIClient, err := paddle.New(cfg.PaddleAPIKey, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		fmt.Printf("Failed to initialize Paddle API client, cause: %s", err)
		return nil
	}

	return app.NewGinWebhookApp(app.GinWebhookAppDependencies{
		PgDB:            pgDB,
		RedisClient:     redisClient,
		Config:          cfg,
		PaddleAPIClient: *paddleAPIClient,
	})
}
