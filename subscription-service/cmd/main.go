package main

import (
	"context"
	"fmt"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/redis"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/compendium-tech/compendium/subscription-service/internal/app"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/joho/godotenv"
)

func main() {
	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	cfg := config.LoadAppConfig()

	tokenManager, err := auth.NewJwtBasedTokenManager(cfg.JwtSingingKey)
	if err != nil {
		fmt.Printf("Failed to initialize token manager, cause: %s", err)
		return
	}

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s", err)
		return
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		fmt.Printf("Failed to connect to Redis, cause: %s", err)
		return
	}

	userService, err := interop.NewGrpcUserServiceClient(cfg.GrpcUserServiceTarget)
	if err != nil {
		fmt.Printf("Failed to initialize user service, cause: %s", err)
		return
	}

	paddleApiClient, err := paddle.New(cfg.PaddleApiKey)
	if err != nil {
		fmt.Printf("Failed to initialize Paddle API client, cause: %s", err)
		return
	}

	app.NewApp(app.Dependencies{
		PgDB:            pgDB,
		RedisClient:     redisClient,
		Config:          cfg,
		TokenManager:    tokenManager,
		UserService:     userService,
		PaddleApiClient: *paddleApiClient,
	}).Run()
}
