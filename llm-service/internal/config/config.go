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
	Environment  string
	GeminiApiKey string
	GrpcPort     uint16
}

func LoadAppConfig() *AppConfig {
	appConfig := &AppConfig{
		Environment:  EnvironmentProd,
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
	}

	env := os.Getenv("ENVIRONMENT")

	if env == "dev" {
		appConfig.Environment = EnvironmentDev
	}

	if grpcPortString := os.Getenv("GRPC_PORT"); grpcPortString != "" {
		var grpcPort uint16
		_, err := fmt.Sscan(grpcPortString, &grpcPort)

		if err == nil {
			appConfig.GrpcPort = grpcPort
		} else {
			log.Printf("Failed to parse gRPC port: %s", grpcPortString)
		}
	}

	return appConfig
}
