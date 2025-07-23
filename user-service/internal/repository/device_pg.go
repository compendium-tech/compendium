package repository

import (
	"context"
	"database/sql"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ztrue/tracerr"
)

const maxDevicesPerUser = 10

type pgDeviceRepository struct {
	db *sql.DB
}

func NewPgDeviceRepository(db *sql.DB) DeviceRepository {
	return &pgDeviceRepository{db: db}
}

func (r *pgDeviceRepository) CreateDevice(ctx context.Context, device model.Device) (finalErr error) {
	tx, err := r.db.Begin()
	if err != nil {
		return tracerr.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// Check the number of devices for the user
	var deviceCount int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE user_id = $1", device.UserID).Scan(&deviceCount)
	if err != nil {
		return tracerr.Errorf("failed to query device count: %w", err)
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
			return tracerr.Errorf("failed to delete oldest device: %w", err)
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
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" { // 23505 is the SQLSTATE for unique_violation
			// Device already exists, fine!
			return nil
		}

		return tracerr.Errorf("failed to insert device: %w", err)
	}

	return tx.Commit()
}

func (r *pgDeviceRepository) DeviceExists(ctx context.Context, userID uuid.UUID, userAgent string, ipAddress string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1 AND user_agent = $2 AND ip_address = $3)`,
		userID, userAgent, ipAddress,
	).Scan(&exists)

	if err != nil {
		return false, tracerr.Wrap(err)
	}

	return exists, nil
}

func (r *pgDeviceRepository) GetDevicesByUserID(ctx context.Context, userID uuid.UUID) ([]model.Device, error) {
	var devices []model.Device

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT * FROM devices WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer rows.Close()

	for rows.Next() {
		var device model.Device
		if err := rows.Scan(&device.ID, &device.UserID, &device.UserAgent, &device.IPAddress, &device.CreatedAt); err != nil {
			return nil, tracerr.Wrap(err)
		}

		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return devices, nil
}

func (r *pgDeviceRepository) RemoveAllDevicesByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM devices WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
