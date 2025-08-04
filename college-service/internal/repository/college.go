package repository

import (
	"context"

	"github.com/compendium-tech/compendium/college-service/internal/model"
)

type CollegeRepository interface {
	SearchColleges(ctx context.Context, semanticSearchText,
		stateOrCountry string, pageIndex, pageSize int) ([]model.College, error)
}
