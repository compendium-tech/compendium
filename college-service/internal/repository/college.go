package repository

import (
	"context"

	"github.com/compendium-tech/compendium/college-service/internal/model"
)

type CollegeRepository interface {
	SearchColleges(ctx context.Context, queryText, stateOrCountry string) ([]model.College, error)
}
