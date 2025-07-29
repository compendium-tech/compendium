package repository

import (
	"context"

	"github.com/compendium-tech/compendium/application-assistant-service/internal/model"
	"github.com/google/uuid"
)

type ApplicationRepository interface {
	CreateApplication(ctx context.Context, app model.Application) error
	GetApplication(ctx context.Context, id uuid.UUID) (*model.Application, error)
	RemoveApplication(ctx context.Context, id uuid.UUID) error
	FindApplicationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Application, error)

	CreateActivity(ctx context.Context, activity model.Activity) error
	GetActivity(ctx context.Context, applicationID uuid.UUID) (*model.Activity, error)
	GetActivities(ctx context.Context, applicationID uuid.UUID) ([]model.Activity, error)
	UpdateActivity(ctx context.Context, activity model.Activity) (*model.Activity, error)
	DeleteActivity(ctx context.Context, applicationID uuid.UUID) error

	CreateHonor(ctx context.Context, honor model.Honor) error
	GetHonor(ctx context.Context, applicationID uuid.UUID) (*model.Honor, error)
	UpdateHonor(ctx context.Context, honor model.Honor) (*model.Honor, error)
	DeleteHonor(ctx context.Context, applicationID uuid.UUID) error
	GetHonors(ctx context.Context, applicationID uuid.UUID) ([]model.Honor, error)
}
