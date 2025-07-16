package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/seacite-tech/compendium/user-service/internal/model"
)

type DeviceRepository interface {
	// If amount of devices exceeds 10, oldest one is removed from database.
	CreateDevice(ctx context.Context, device model.Device) error
	DeviceExists(userId uuid.UUID, userAgent string, ipAddress string) (bool, error)
}
