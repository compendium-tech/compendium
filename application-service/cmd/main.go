package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/compendium-tech/compendium/application-service/internal/app"
	"github.com/compendium-tech/compendium/application-service/internal/config"
	"github.com/compendium-tech/compendium/application-service/internal/interop"
)

func main() {
	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	cfg := config.LoadAppConfig()

	tokenManager, err := auth.NewJwtBasedTokenManager(cfg.JwtSingingKey)
	if err != nil {
		fmt.Printf("Failed to initialize token manager, cause: %v\n", err)
		return
	}

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %v\n", err)
		return
	}

	llmService, err := interop.NewGrpcLLMServiceClient(cfg.GrpcLLMServiceClientTarget)
	if err != nil {
		fmt.Printf("Failed to initialize llm service client, cause: %v\n", err)
		return
	}

	deps := app.Dependencies{
		Config:       cfg,
		TokenManager: tokenManager,
		PgDB:         pgDB,
		LLMService:   llmService,
	}

	err = app.NewApp(deps).Run()
	if err != nil {
		fmt.Printf("Failed to run application service, cause: %v\n", err)
	}
}
