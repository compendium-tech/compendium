package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/redis"
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/compendium-tech/compendium/user-service/internal/app"
	"github.com/compendium-tech/compendium/user-service/internal/config"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	"github.com/compendium-tech/compendium/user-service/internal/geoip"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	"github.com/compendium-tech/compendium/user-service/internal/ua"
)

func main() {
	appMode := flag.String("mode", "", "Specify the application mode: 'http' for Gin app or 'grpc' for gRPC app")
	flag.Parse()

	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	switch *appMode {
	case "http":
		runHttpApp(ctx)
	case "grpc":
		runGrpcApp(ctx)
	default:
		fmt.Printf("Invalid application mode specified: %s. Please use 'http' or 'grpc'.\n", *appMode)
	}
}

func runHttpApp(ctx context.Context) {
	fmt.Println("Starting Gin (HTTP) application...")

	cfg := config.LoadGinAppConfig()

	tokenManager, err := auth.NewJwtBasedTokenManager(cfg.JwtSingingKey)
	if err != nil {
		fmt.Printf("Failed to initialize token manager, cause: %s\n", err)
		return
	}

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s\n", err)
		return
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		fmt.Printf("Failed to connect to Redis, cause: %s\n", err)
		return
	}

	kafkaEmailSender := email.NewKafkaEmailMessageProducer(cfg.EmailDeliveryKafkaBroker, cfg.EmailDeliveryKafkaTopic)

	emailMessageBuilder, err := email.NewMessageBuilder()
	if err != nil {
		fmt.Printf("Failed to initialize email builder, cause: %s\n", err)
		return
	}

	geoIP := geoip.NewGeoIP2Client(cfg.GeoIP2AccountID, cfg.GeoIP2LicenseKey, cfg.GeoIP2Host)

	userAgentParser := ua.NewUserAgentParser()

	deps := app.GinAppDependencies{
		PgDB:                pgDB,
		RedisClient:         redisClient,
		Config:              cfg,
		TokenManager:        tokenManager,
		EmailSender:         kafkaEmailSender,
		EmailMessageBuilder: emailMessageBuilder,
		GeoIP:               geoIP,
		UserAgentParser:     userAgentParser,
		PasswordHasher:      hash.NewBcryptPasswordHasher(bcrypt.DefaultCost),
	}

	_ = app.NewGinApp(deps).Run()
}

func runGrpcApp(ctx context.Context) {
	fmt.Println("Starting gRPC application...")

	cfg := config.LoadGrpcAppConfig()

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s\n", err)
		return
	}

	deps := app.GrpcAppDependencies{
		PgDB:   pgDB,
		Config: cfg,
	}

	_ = app.NewGrpcApp(deps).Run()
}
