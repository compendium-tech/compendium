package pg

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunUpPgMigrations(ctx context.Context, db *sql.DB, migrationsDir string) error {
	log.Println("Running UP PostgreSQL migrations...")
	return runMigrations(ctx, db, migrationsDir, ".up.sql", false)
}

func RunDownPgMigrations(ctx context.Context, db *sql.DB, migrationsDir string) error {
	log.Println("Running DOWN PostgreSQL migrations...")
	return runMigrations(ctx, db, migrationsDir, ".down.sql", true)
}

func runMigrations(ctx context.Context, db *sql.DB, migrationsDir string, suffix string, reverse bool) error {
	log.Printf("Running migrations with suffix '%s' in directory '%s'...", suffix, migrationsDir)

	// Read all entries in the migrations directory.
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory %s: %v", migrationsDir, err)
	}

	var migrationFiles []fs.DirEntry
	// Filter files to include only those ending with the specified suffix.
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), suffix) {
			migrationFiles = append(migrationFiles, file)
		}
	}

	// Sort the migration files alphabetically by name to ensure consistent order.
	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Name() < migrationFiles[j].Name()
	})

	// If 'reverse' is true, reverse the order of the sorted files.
	if reverse {
		for i, j := 0, len(migrationFiles)-1; i < j; i, j = i+1, j-1 {
			migrationFiles[i], migrationFiles[j] = migrationFiles[j], migrationFiles[i]
		}
	}

	// Iterate through the selected migration files and execute their SQL content.
	for _, file := range migrationFiles {
		filePath := filepath.Join(migrationsDir, file.Name())
		log.Printf("Executing migration: %s", filePath)

		// Read the content of the SQL file.
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
		}

		// Execute the SQL content against the database.
		if _, err := db.ExecContext(ctx, string(sqlContent)); err != nil {
			return fmt.Errorf("failed to execute SQL from %s: %w", filePath, err)
		}
	}

	log.Printf("Migrations with suffix '%s' complete.", suffix)
	return nil
}
