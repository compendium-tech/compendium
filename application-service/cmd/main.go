package cmd

import (
	"context"
	"fmt"
	"github.com/compendium-tech/compendium/application-service/internal/app"
	"github.com/compendium-tech/compendium/application-service/internal/config"
	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	llmservice "github.com/compendium-tech/compendium/llm-common/pkg/service"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
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
		fmt.Printf("Failed to initialize token manager, cause: %s\n", err)
		return
	}

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s\n", err)
		return
	}

	geminiAPI, err := llmservice.NewGeminiClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		fmt.Printf("Failed to initialize gemini client, cause: %s\n", err)
		return
	}

	deps := app.Dependencies{
		Config:       cfg,
		TokenManager: tokenManager,
		PgDB:         pgDB,
		LLMService:   geminiAPI,
	}

	_ = app.NewApp(deps).Run()
}
