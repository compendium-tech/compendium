package domain

import (
	"time"

	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/google/uuid"
)

type ApplicationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type ActivityResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Role         string                 `json:"role"`
	Description  *string                `json:"description"`
	HoursPerWeek int                    `json:"hoursPerWeek"`
	WeeksPerYear int                    `json:"weeksPerYear"`
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
	HoursPerWeek int                    `json:"hoursPerWeek"`
	WeeksPerYear int                    `json:"weeksPerYear"`
	Category     model.ActivityCategory `json:"category"`
	Grades       []model.Grade          `json:"grades"`
}

type UpdateHonorRequest struct {
	Title       string           `json:"title"`
	Description *string          `json:"description"`
	Level       model.HonorLevel `json:"level"`
	Grade       model.Grade      `json:"grade"`
}

type CreateEssayRequest struct {
	Kind    model.EssayType `json:"type"`
	Content string          `json:"content"`
}

type UpdateEssayRequest struct {
	Kind    model.EssayType `json:"type"`
	Content string          `json:"content"`
}

type EssayResponse struct {
	ID      uuid.UUID       `json:"id"`
	Kind    model.EssayType `json:"type"`
	Content string          `json:"content"`
}

type CreateSupplementalEssayRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateSupplementalEssayRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SupplementalEssayResponse struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}
