package main

import (
	"context"
	"fmt"

	"github.com/compendium-tech/compendium/common/pkg/pg"

	"github.com/compendium-tech/compendium/user-service/internal/app"
	"github.com/compendium-tech/compendium/user-service/internal/config"
)

func runGrpcApp(ctx context.Context) {
	fmt.Println("Starting gRPC application...")

	cfg := config.LoadGrpcAppConfig()

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL, cause: %s\n", err)
		return
	}

	deps := app.GrpcAppDependencies{
		PgDB:   pgDB,
		Config: cfg,
	}

	_ = app.NewGrpcApp(deps).Run()
}
