package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPgClient(ctx context.Context, host string, port uint16, username, password, databaseName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, databaseName,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
