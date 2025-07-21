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
	Environment              string
	PgHost                   string
	PgPort                   uint16
	PgUsername               string
	PgPassword               string
	PgDatabaseName           string
	RedisHost                string
	RedisPort                uint16
	JwtSingingKey            string
	CsrfTokenHashSalt        string
	EmailDeliveryKafkaBroker string
	EmailDeliveryKafkaTopic  string
}

func LoadAppConfig() AppConfig {
	appConfig := AppConfig{
		PgHost:                   os.Getenv("POSTGRES_HOST"),
		PgUsername:               os.Getenv("POSTGRES_USERNAME"),
		PgPassword:               os.Getenv("POSTGRES_PASSWORD"),
		PgDatabaseName:           os.Getenv("POSTGRES_DATABASE_NAME"),
		RedisHost:                os.Getenv("REDIS_HOST"),
		JwtSingingKey:            os.Getenv("JWT_SIGNING_KEY"),
		CsrfTokenHashSalt:        os.Getenv("CSRF_TOKEN_HASH_SALT"),
		EmailDeliveryKafkaBroker: os.Getenv("EMAIL_DELIVERY_KAFKA_BROKER"),
		EmailDeliveryKafkaTopic:  os.Getenv("EMAIL_DELIVERY_KAFKA_TOPIC"),
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
