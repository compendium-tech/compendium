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
	ElasticsearchHost          string
	ElasticsearchPort          uint16
	ElasticsearchUsername      string
	ElasticsearchPassword      string
	JwtSingingKey              string
	GrpcLLMServiceClientTarget string
	CsrfTokenHashSalt          string
}

func LoadAppConfig() *AppConfig {
	appConfig := &AppConfig{
		Environment:                EnvironmentProd,
		ElasticsearchHost:          os.Getenv("ELASTICSEARCH_HOST"),
		ElasticsearchUsername:      os.Getenv("ELASTICSEARCH_USERNAME"),
		ElasticsearchPassword:      os.Getenv("ELASTICSEARCH_PASSWORD"),
		JwtSingingKey:              os.Getenv("JWT_SIGNING_KEY"),
		CsrfTokenHashSalt:          os.Getenv("CSRF_TOKEN_HASH_SALT"),
		GrpcLLMServiceClientTarget: os.Getenv("GRPC_LLM_SERVICE_CLIENT_TARGET"),
	}

	env := os.Getenv("ENVIRONMENT")

	if env == "dev" {
		appConfig.Environment = EnvironmentDev
	}

	if port := os.Getenv("ELASTICSEARCH_PORT"); port != "" {
		var esPort uint16
		_, err := fmt.Sscan(port, &esPort)

		if err == nil {
			appConfig.ElasticsearchPort = esPort
		} else {
			log.Printf("Failed to parse postgres port: %s", port)
		}
	}

	return appConfig
}
