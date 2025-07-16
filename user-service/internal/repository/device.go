package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/seacite-tech/compendium/user-service/internal/model"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device model.Device) error
	DeviceExists(userId uuid.UUID, userAgent string, ipAddress string) (bool, error)
}
