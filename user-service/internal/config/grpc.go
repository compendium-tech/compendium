package config

import (
	"fmt"
	"log"
	"os"
)

type GrpcAppConfig struct {
	Environment    string
	PgHost         string
	PgPort         uint16
	PgUsername     string
	PgPassword     string
	PgDatabaseName string
}

func LoadGrpcAppConfig() GrpcAppConfig {
	appConfig := GrpcAppConfig{
		PgHost:         os.Getenv("POSTGRES_HOST"),
		PgUsername:     os.Getenv("POSTGRES_USERNAME"),
		PgPassword:     os.Getenv("POSTGRES_PASSWORD"),
		PgDatabaseName: os.Getenv("POSTGRES_DATABASE_NAME"),
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
