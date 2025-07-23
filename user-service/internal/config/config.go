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

const (
	ModeHttp string = "HTTP"
	ModeGrpc string = "GRPC"
)

type AppConfig struct {
	Mode                     string
	Environment              string
	GrpcPort                 uint16
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
	GeoIP2AccountID          string
	GeoIP2LicenseKey         string
	GeoIP2Host               string
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
		GeoIP2AccountID:          os.Getenv("GEOIP2_ACCOUNT_ID"),
		GeoIP2LicenseKey:         os.Getenv("GEOIP2_LICENSE_KEY"),
		GeoIP2Host:               os.Getenv("GEOIP2_HOST"),
	}

	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		appConfig.Environment = EnvironmentDev
	case "prod":
		appConfig.Environment = EnvironmentProd
	default:
		log.Printf("Unknown `environment` variable: %s, defaulting to dev", env)
		appConfig.Environment = EnvironmentDev
	}

	mode := os.Getenv("MODE")
	switch mode {
	case "http":
		appConfig.Mode = ModeHttp
	case "grpc":
		appConfig.Mode = ModeGrpc
	default:
		log.Printf("Unknown `mode` value: %s, defaulting to http", mode)
		appConfig.Mode = ModeHttp
	}

	if port := os.Getenv("GRPC_PORT"); port != "" {
		var grpcPort uint16
		_, err := fmt.Sscan(port, &grpcPort)

		if err == nil {
			appConfig.GrpcPort = grpcPort
		} else {
			log.Printf("Failed to parse GRPC port: %s", port)
		}
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
