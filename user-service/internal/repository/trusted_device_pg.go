package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/user-service/internal/model"
)

const maxDevicesPerUser = 10

type pgTrustedDeviceRepository struct {
	db *sql.DB
}

func NewPgTrustedDeviceRepository(db *sql.DB) TrustedDeviceRepository {
	return &pgTrustedDeviceRepository{db: db}
}

func (r *pgTrustedDeviceRepository) UpsertDevice(ctx context.Context, device model.TrustedDevice) {
	tx, err := r.db.Begin()
	if err != nil {
		panic(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback()

	// Check if the device already exists
	var exists bool
	err = tx.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1 AND user_agent = $2 AND ip_address = $3)`,
		device.UserID, device.UserAgent, device.IPAddress,
	).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if exists {
		return
	}

	// Check the number of devices for the user
	var deviceCount int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE user_id = $1", device.UserID).Scan(&deviceCount)
	if err != nil {
		panic(fmt.Errorf("failed to query device count: %w", err))
	}

	// If amount of devices exceeds a limit, oldest one is removed from database.
	if deviceCount >= maxDevicesPerUser {
		_, err := tx.ExecContext(
			ctx,
			`DELETE FROM devices WHERE id = (
				SELECT id FROM devices WHERE user_id = $1 ORDER BY created_at ASC LIMIT 1
			)`,
			device.UserID,
		)

		if err != nil {
			panic(fmt.Errorf("failed to delete oldest device: %w", err))
		}
	}

	// Insert the new device
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO devices (id, user_id, user_agent, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		device.ID, device.UserID, device.UserAgent, device.IPAddress, device.CreatedAt,
	)
	if err != nil {
		panic(fmt.Errorf("failed to insert device: %w", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("failed to commit transaction: %w", err))
	}
}

func (r *pgTrustedDeviceRepository) DeviceExists(ctx context.Context, userID uuid.UUID, userAgent string, ipAddress string) bool {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1 AND user_agent = $2 AND ip_address = $3)`,
		userID, userAgent, ipAddress,
	).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
