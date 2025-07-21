package main

import (
	"context"
	"fmt"

	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/redis"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	emailDelivery "github.com/compendium-tech/compendium/email-delivery-service/pkg/email"
	"github.com/compendium-tech/compendium/user-service/internal/app"
	"github.com/compendium-tech/compendium/user-service/internal/config"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	"github.com/compendium-tech/compendium/user-service/pkg/auth"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
		return
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

	kafkaEmailSender := emailDelivery.NewKafkaEmailMessageProducer(cfg.EmailDeliveryKafkaBroker, cfg.EmailDeliveryKafkaTopic)
	emailMessageBuilder, err := email.NewEmailMessageBuilder()
	if err != nil {
		fmt.Printf("Failed to initialize email builder, cause: %s", err)
		return
	}

	app.NewApp(app.Dependencies{
		PgDb:                pgDB,
		RedisClient:         redisClient,
		Config:              cfg,
		TokenManager:        tokenManager,
		EmailSender:         kafkaEmailSender,
		EmailMessageBuilder: emailMessageBuilder,
		PasswordHasher:      hash.NewBcryptPasswordHasher(bcrypt.DefaultCost),
	}).Run()
}
