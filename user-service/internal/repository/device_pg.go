package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seacite-tech/compendium/user-service/internal/model"
	"github.com/ztrue/tracerr"
)

type PostgresDeviceRepository struct {
	db *sql.DB
}

func NewPgDeviceRepository(db *sql.DB) *PostgresDeviceRepository {
	return &PostgresDeviceRepository{db: db}
}

func (r *PostgresDeviceRepository) CreateDevice(ctx context.Context, device model.Device) error {
	tx, err := r.db.Begin()
	if err != nil {
		return tracerr.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check the number of devices for the user
	var deviceCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM devices WHERE user_id = $1", device.UserId).Scan(&deviceCount)
	if err != nil {
		return tracerr.Errorf("failed to query device count: %w", err)
	}

	// If amount of devices exceeds 10, oldest one is removed from database.
	if deviceCount >= 10 {
		_, err := tx.Exec(
			`DELETE FROM devices WHERE id = (
				SELECT id FROM devices WHERE user_id = $1 ORDER BY created_at ASC LIMIT 1
			)`,
			device.UserId,
		)
		if err != nil {
			return tracerr.Errorf("failed to delete oldest device: %w", err)
		}
	}

	// Insert the new device
	_, err = tx.Exec(
		`INSERT INTO devices (id, user_id, user_agent, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		device.Id, device.UserId, device.UserAgent, device.IpAddress, device.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
			// Device already exists, fine!
			return nil
		}
		return tracerr.Errorf("failed to insert device: %w", err)
	}

	return tx.Commit()
}

func (r *PostgresDeviceRepository) DeviceExists(userId uuid.UUID, userAgent string, ipAddress string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1 AND user_agent = $2 AND ip_address = $3)`,
		userId, userAgent, ipAddress,
	).Scan(&exists)

	if err != nil {
		return false, tracerr.Wrap(err)
	}

	return exists, nil
}
