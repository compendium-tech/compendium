package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"github.com/compendium-tech/compendium/llm-service/internal/app"
	"github.com/compendium-tech/compendium/llm-service/internal/config"
	"github.com/compendium-tech/compendium/llm-service/internal/service"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	cfg := config.LoadAppConfig()

	geminiAPIClient, err := service.NewGeminiClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		fmt.Printf("Failed to initialize Gemini API client, cause: %s", err)
		return
	}

	err = app.NewApp(app.Dependencies{
		Config:     cfg,
		LLMService: geminiAPIClient,
	}).Run()
	if err != nil {
		fmt.Printf("Failed to start LLM service, cause: %v\n", err)
	}
}
