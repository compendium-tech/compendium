package repository

import (
	"context"

	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/google/uuid"
)

type ApplicationRepository interface {
	GetApplication(ctx context.Context, id uuid.UUID) *model.Application
	FindApplicationsByUserID(ctx context.Context, userID uuid.UUID) []model.Application

	CreateApplication(ctx context.Context, app model.Application)
	UpdateApplicationName(ctx context.Context, applicationID uuid.UUID, name string)
	RemoveApplication(ctx context.Context, id uuid.UUID)

	GetActivities(ctx context.Context, applicationID uuid.UUID) []model.Activity
	PutActivities(ctx context.Context, applicationID uuid.UUID, activities []model.Activity)

	GetHonors(ctx context.Context, applicationID uuid.UUID) []model.Honor
	PutHonors(ctx context.Context, applicationID uuid.UUID, honors []model.Honor)

	GetEssays(ctx context.Context, applicationID uuid.UUID) []model.Essay
	PutEssays(ctx context.Context, applicationID uuid.UUID, essays []model.Essay)

	GetSupplementalEssays(ctx context.Context, applicationID uuid.UUID) []model.SupplementalEssay
	PutSupplementalEssays(ctx context.Context, applicationID uuid.UUID, essays []model.SupplementalEssay)
}
