package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device model.Device) error
	DeviceExists(ctx context.Context, userID uuid.UUID, userAgent string, ipAddress string) (bool, error)
	GetDevicesByUserID(ctx context.Context, userID uuid.UUID) ([]model.Device, error)
	RemoveAllDevicesByUserID(ctx context.Context, userID uuid.UUID) error
}
