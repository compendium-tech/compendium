package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type TrustedDeviceRepository interface {
	ExistsOrCreateDevice(ctx context.Context, device model.TrustedDevice) error
	DeviceExists(ctx context.Context, userID uuid.UUID, userAgent string, ipAddress string) (bool, error)
}
