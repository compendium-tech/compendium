package domain

import (
	"time"

	"github.com/compendium-tech/compendium/application-assistant-service/internal/model"
	"github.com/google/uuid"
)

type ApplicationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ActivityResponse struct {
	Name         string                 `json:"name"`
	Role         string                 `json:"role"`
	Description  *string                `json:"description"`
	HoursPerWeek int                    `json:"hours_per_week"`
	WeeksPerYear int                    `json:"weeks_per_year"`
	Category     model.ActivityCategory `json:"category"`
	Grades       []model.Grade          `json:"grades"`
}

type HonorResponse struct {
	Title       string           `json:"title"`
	Description *string          `json:"description"`
	Level       model.HonorLevel `json:"level"`
	Grade       model.Grade      `json:"grade"`
}

type CreateApplicationRequest struct {
	Name string `json:"name"`
}

type UpdateActivityRequest struct {
	Name         string                 `json:"name"`
	Role         string                 `json:"role"`
	Description  *string                `json:"description"`
	HoursPerWeek int                    `json:"hours_per_week"`
	WeeksPerYear int                    `json:"weeks_per_year"`
	Category     model.ActivityCategory `json:"category"`
	Grades       []model.Grade          `json:"grades"`
}

type UpdateHonorRequest struct {
	Title       string           `json:"title"`
	Description *string          `json:"description"`
	Level       model.HonorLevel `json:"level"`
	Grade       model.Grade      `json:"grade"`
}
