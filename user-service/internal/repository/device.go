package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device model.Device) error
	DeviceExists(ctx context.Context, userId uuid.UUID, userAgent string, ipAddress string) (bool, error)
	GetDevicesByUserId(ctx context.Context, userId uuid.UUID) ([]model.Device, error)
	RemoveAllDevicesByUserId(ctx context.Context, userId uuid.UUID) error
}
