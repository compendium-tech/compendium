package config

import (
	"fmt"
	"log"
	"os"
)

const (
	EnvironmentDev  string = "dev"
	EnvironmentProd string = "prod"
)

type ProductIDs struct {
	StudentSubscriptionProductID   string
	TeamSubscriptionProductID      string
	CommunitySubscriptionProductID string
}

type AppConfig struct {
	Environment                 string
	PgHost                      string
	PgPort                      uint16
	PgUsername                  string
	PgPassword                  string
	PgDatabaseName              string
	RedisHost                   string
	RedisPort                   uint16
	JwtSingingKey               string
	CsrfTokenHashSalt           string
	PaddleWebhookSecret         string
	ProductIDs                  ProductIDs
	PaddleAPIKey                string
	GrpcUserServiceClientTarget string
}

func LoadAppConfig() *AppConfig {
	appConfig := &AppConfig{
		Environment:         EnvironmentProd,
		PgHost:              os.Getenv("POSTGRES_HOST"),
		PgUsername:          os.Getenv("POSTGRES_USERNAME"),
		PgPassword:          os.Getenv("POSTGRES_PASSWORD"),
		PgDatabaseName:      os.Getenv("POSTGRES_DATABASE_NAME"),
		RedisHost:           os.Getenv("REDIS_HOST"),
		JwtSingingKey:       os.Getenv("JWT_SIGNING_KEY"),
		CsrfTokenHashSalt:   os.Getenv("CSRF_TOKEN_HASH_SALT"),
		PaddleWebhookSecret: os.Getenv("PADDLE_WEBHOOK_SECRET"),
		ProductIDs: ProductIDs{
			StudentSubscriptionProductID:   os.Getenv("STUDENT_SUBSCRIPTION_PRODUCT_ID"),
			TeamSubscriptionProductID:      os.Getenv("TEAM_SUBSCRIPTION_PRODUCT_ID"),
			CommunitySubscriptionProductID: os.Getenv("COMMUNITY_SUBSCRIPTION_PRODUCT_ID"),
		},
		PaddleAPIKey:                os.Getenv("PADDLE_API_KEY"),
		GrpcUserServiceClientTarget: os.Getenv("GRPC_USER_SERVICE_CLIENT_TARGET"),
	}

	env := os.Getenv("ENVIRONMENT")

	if env == "dev" {
		appConfig.Environment = EnvironmentDev
	}

	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		var pgPort uint16
		_, err := fmt.Sscan(port, &pgPort)

		if err == nil {
			appConfig.PgPort = pgPort
		} else {
			log.Printf("Failed to parse postgres port: %s", port)
		}
	}

	if port := os.Getenv("REDIS_PORT"); port != "" {
		var redisPort uint16
		_, err := fmt.Sscan(port, &redisPort)

		if err == nil {
			appConfig.RedisPort = redisPort
		} else {
			log.Printf("Failed to parse redis port: %s", port)
		}
	}

	return appConfig
}
