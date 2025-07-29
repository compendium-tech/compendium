package repository

import (
	"context"

	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/google/uuid"
)

type ApplicationRepository interface {
	GetApplication(ctx context.Context, id uuid.UUID) (*model.Application, error)
	FindApplicationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Application, error)

	CreateApplication(ctx context.Context, app model.Application) error
	UpdateApplicationName(ctx context.Context, applicationID uuid.UUID, name string) error
	RemoveApplication(ctx context.Context, id uuid.UUID) error

	GetActivities(ctx context.Context, applicationID uuid.UUID) ([]model.Activity, error)
	PutActivities(ctx context.Context, applicationID uuid.UUID, activities []model.Activity) error

	GetHonors(ctx context.Context, applicationID uuid.UUID) ([]model.Honor, error)
	PutHonors(ctx context.Context, applicationID uuid.UUID, honors []model.Honor) error

	GetEssays(ctx context.Context, applicationID uuid.UUID) ([]model.Essay, error)
	PutEssays(ctx context.Context, applicationID uuid.UUID, essays []model.Essay) error

	GetSupplementalEssays(ctx context.Context, applicationID uuid.UUID) ([]model.SupplementalEssay, error)
	PutSupplementalEssays(ctx context.Context, applicationID uuid.UUID, essays []model.SupplementalEssay) error
}
