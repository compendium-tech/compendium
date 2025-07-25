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
		PgHost:              "127.0.0.1",
		PgPort:              5432,
		PgUsername:          "postgres",
		PgPassword:          "",
		PgDatabaseName:      "",
		RedisPort:           6379,
		RedisHost:           "127.0.0.1",
		JwtSingingKey:       "",
		CsrfTokenHashSalt:   "",
		PaddleWebhookSecret: "",
		ProductIDs: ProductIDs{
			StudentSubscriptionProductID:   "",
			TeamSubscriptionProductID:      "",
			CommunitySubscriptionProductID: "",
		},
		PaddleAPIKey:                "",
		GrpcUserServiceClientTarget: "localhost:2000",
	}

	env := os.Getenv("ENVIRONMENT")

	if env == "dev" {
		appConfig.Environment = EnvironmentDev
	}

	appConfig.PgHost = os.Getenv("POSTGRES_HOST")
	appConfig.PgUsername = os.Getenv("POSTGRES_USERNAME")
	appConfig.PgPassword = os.Getenv("POSTGRES_PASSWORD")
	appConfig.PgDatabaseName = os.Getenv("POSTGRES_DATABASE_NAME")
	appConfig.RedisHost = os.Getenv("REDIS_HOST")
	appConfig.JwtSingingKey = os.Getenv("JWT_SIGNING_KEY")
	appConfig.CsrfTokenHashSalt = os.Getenv("CSRF_TOKEN_HASH_SALT")
	appConfig.PaddleWebhookSecret = os.Getenv("PADDLE_WEBHOOK_SECRET")
	appConfig.ProductIDs = ProductIDs{
		StudentSubscriptionProductID:   os.Getenv("STUDENT_SUBSCRIPTION_PRODUCT_ID"),
		TeamSubscriptionProductID:      os.Getenv("TEAM_SUBSCRIPTION_PRODUCT_ID"),
		CommunitySubscriptionProductID: os.Getenv("COMMUNITY_SUBSCRIPTION_PRODUCT_ID"),
	}
	appConfig.PaddleAPIKey = os.Getenv("PADDLE_API_KEY")
	appConfig.GrpcUserServiceClientTarget = os.Getenv("GRPC_USER_SERVICE_CLIENT_TARGET")

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
