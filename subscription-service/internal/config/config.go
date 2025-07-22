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

type AppConfig struct {
	Environment         string
	PgHost              string
	PgPort              uint16
	PgUsername          string
	PgPassword          string
	PgDatabaseName      string
	RedisHost           string
	RedisPort           uint16
	JwtSingingKey       string
	CsrfTokenHashSalt   string
	PaddleWebhookSecret string
	PaddleProductIds    PaddleProductIds
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
		PaddleProductIds: PaddleProductIds{
			StudentSubscriptionProductId:   "",
			TeamSubscriptionProductId:      "",
			CommunitySubscriptionProductId: "",
		},
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
	appConfig.PaddleProductIds = PaddleProductIds{
		StudentSubscriptionProductId:   os.Getenv("PADDLE_STUDENT_SUBSCRIPTION_PRICE_ID"),
		TeamSubscriptionProductId:      os.Getenv("PADDLE_TEAM_SUBSCRIPTION_PRICE_ID"),
		CommunitySubscriptionProductId: os.Getenv("PADDLE_COMMUNITY_SUBSCRIPTION_PRICE_ID"),
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
