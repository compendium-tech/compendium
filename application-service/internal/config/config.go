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
	Environment                string
	PgHost                     string
	PgPort                     uint16
	PgUsername                 string
	PgPassword                 string
	PgDatabaseName             string
	JwtSingingKey              string
	GrpcLLMServiceClientTarget string
	CsrfTokenHashSalt          string
}

func LoadAppConfig() *AppConfig {
	appConfig := &AppConfig{
		Environment:                EnvironmentProd,
		PgHost:                     os.Getenv("POSTGRES_HOST"),
		PgUsername:                 os.Getenv("POSTGRES_USERNAME"),
		PgPassword:                 os.Getenv("POSTGRES_PASSWORD"),
		PgDatabaseName:             os.Getenv("POSTGRES_DATABASE_NAME"),
		JwtSingingKey:              os.Getenv("JWT_SIGNING_KEY"),
		CsrfTokenHashSalt:          os.Getenv("CSRF_TOKEN_HASH_SALT"),
		GrpcLLMServiceClientTarget: os.Getenv("GRPC_LLM_SERVICE_CLIENT_TARGET"),
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

	return appConfig
}
