package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/seacite-tech/compendium/common/pkg/pg"
	"github.com/seacite-tech/compendium/common/pkg/redis"
	"github.com/seacite-tech/compendium/user-service/internal/app"
	"github.com/seacite-tech/compendium/user-service/internal/config"
	"github.com/seacite-tech/compendium/user-service/internal/email"
	"github.com/seacite-tech/compendium/user-service/internal/validate"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
)

func main() {
	validate.InitializeValidator()

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
	}

	smtpEmailSender, err := email.NewSmtpEmailSender(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpFrom)
	if err != nil {
		fmt.Printf("Failed to initialize email service, cause: %s", err)
		return
	}

	app.NewApp(app.Dependencies{
		PgDb:         pgDB,
		RedisClient:  redisClient,
		Config:       cfg,
		TokenManager: tokenManager,
		EmailSender:  smtpEmailSender,
	}).Run()
}
