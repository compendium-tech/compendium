package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/compendium-tech/compendium/common/pkg/validate"
)

func main() {
	appMode := flag.String("mode", "", "Specify the application mode: 'http' for Gin app or 'grpc' for gRPC app")
	flag.Parse()

	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	switch *appMode {
	case "http":
		runHttpApp(ctx)
	case "grpc":
		runGrpcApp(ctx)
	default:
		fmt.Printf("Invalid application mode specified: %s. Please use 'http' or 'grpc'.\n", *appMode)
	}
}
