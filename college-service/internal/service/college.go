package service

import (
	"context"

	"github.com/compendium-tech/compendium/common/pkg/log"

	"github.com/compendium-tech/compendium/college-service/internal/domain"
	"github.com/compendium-tech/compendium/college-service/internal/repository"
)

type CollegeService interface {
	SearchColleges(ctx context.Context, request domain.SearchCollegesRequest) ([]domain.CollegeResponse, error)
}

type collegeService struct {
	collegeRepository repository.CollegeRepository
}

func NewCollegeService(collegeRepository repository.CollegeRepository) CollegeService {
	return &collegeService{
		collegeRepository: collegeRepository,
	}
}

func (a *collegeService) SearchColleges(ctx context.Context, request domain.SearchCollegesRequest) ([]domain.CollegeResponse, error) {
	logger := log.L(ctx).WithField("request", request)
	logger.Info("Searching for colleges")

	pageIndex := 0
	if request.PageIndex != nil {
		pageIndex = *request.PageIndex
	}

	pageSize := 10

	semanticSearchText := ""
	if request.SemanticSearchText != nil {
		semanticSearchText = *request.SemanticSearchText
	}

	stateOrCountry := ""
	if request.StateOrCountry != nil {
		stateOrCountry = *request.StateOrCountry
	}

	colleges, err := a.collegeRepository.SearchColleges(
		ctx,
		semanticSearchText,
		stateOrCountry,
		pageIndex,
		pageSize,
	)
	if err != nil {
		return nil, err
	}

	var collegeResponses []domain.CollegeResponse
	for _, college := range colleges {
		collegeResponses = append(collegeResponses, domain.CollegeResponse{
			Name:           college.Name,
			City:           college.City,
			StateOrCountry: college.StateOrCountry,
			Description:    college.Description,
		})
	}

	logger.Infof("Found %d colleges for search", len(collegeResponses))
	return collegeResponses, nil
}
