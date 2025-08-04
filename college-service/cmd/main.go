package main

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/joho/godotenv"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/compendium-tech/compendium/college-service/internal/app"
	"github.com/compendium-tech/compendium/college-service/internal/config"
)

func main() {
	validate.InitValidator()

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

	elasticsearchClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:  cfg.ElasticsearchUsername,
		Password:  cfg.ElasticsearchPassword,
		Addresses: []string{fmt.Sprintf("%s:%d", cfg.ElasticsearchHost, cfg.ElasticsearchPort)},
	})
	if err != nil {
		fmt.Printf("Failed to initialize elasticsearch client, cause: %v\n", err)
		return
	}

	deps := app.Dependencies{
		Config:              cfg,
		TokenManager:        tokenManager,
		ElasticsearchClient: elasticsearchClient,
	}

	err = app.NewApp(deps).Run()
	if err != nil {
		fmt.Printf("Failed to run college service, cause: %v\n", err)
	}
}
