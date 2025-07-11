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
	Environment       string
	PgHost            string
	PgPort            uint16
	PgUsername        string
	PgPassword        string
	PgDatabaseName    string
	RedisHost         string
	RedisPort         uint16
	SmtpHost          string
	SmtpPort          uint16
	SmtpUsername      string
	SmtpPassword      string
	SmtpFrom          string
	JwtSingingKey     string
	CsrfTokenHashSalt string
}

func LoadAppConfig() *AppConfig {
	appConfig := &AppConfig{
		Environment:       EnvironmentProd,
		PgHost:            "127.0.0.1",
		PgPort:            5432,
		PgUsername:        "postgres",
		PgPassword:        "",
		PgDatabaseName:    "",
		RedisPort:         6379,
		RedisHost:         "127.0.0.1",
		SmtpHost:          "",
		SmtpPort:          534,
		SmtpPassword:      "",
		SmtpFrom:          "",
		JwtSingingKey:     "",
		CsrfTokenHashSalt: "",
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
	appConfig.SmtpHost = os.Getenv("SMTP_HOST")
	appConfig.JwtSingingKey = os.Getenv("JWT_SIGNING_KEY")
	appConfig.CsrfTokenHashSalt = os.Getenv("CSRF_TOKEN_HASH_SALT")

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

	if port := os.Getenv("SMTP_PORT"); port != "" {
		var smtpPort uint16
		_, err := fmt.Sscan(port, &smtpPort)

		if err == nil {
			appConfig.SmtpPort = smtpPort
		} else {
			log.Printf("Failed to parse smtp port: %s", port)
		}
	}

	appConfig.SmtpUsername = os.Getenv("SMTP_USERNAME")
	appConfig.SmtpPassword = os.Getenv("SMTP_PASSWORD")
	appConfig.SmtpFrom = os.Getenv("SMTP_FROM")

	return appConfig
}
