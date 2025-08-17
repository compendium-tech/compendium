package service

import (
	"context"

	"github.com/compendium-tech/compendium/common/pkg/log"

	"github.com/compendium-tech/compendium/college-service/internal/domain"
	"github.com/compendium-tech/compendium/college-service/internal/repository"
)

type CollegeService interface {
	SearchColleges(ctx context.Context, request domain.SearchCollegesRequest) []domain.CollegeResponse
}

type collegeService struct {
	collegeRepository repository.CollegeRepository
}

func NewCollegeService(collegeRepository repository.CollegeRepository) CollegeService {
	return &collegeService{
		collegeRepository: collegeRepository,
	}
}

func (a *collegeService) SearchColleges(ctx context.Context, request domain.SearchCollegesRequest) []domain.CollegeResponse {
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

	colleges := a.collegeRepository.SearchColleges(
		ctx,
		semanticSearchText,
		stateOrCountry,
		pageIndex,
		pageSize,
	)

	collegesResponse := make([]domain.CollegeResponse, len(colleges))
	for i, college := range colleges {
		collegesResponse[i] = domain.CollegeResponse{
			Name:           college.Name,
			City:           college.City,
			StateOrCountry: college.StateOrCountry,
			Description:    college.Description,
		}
	}

	logger.Infof("Found %d colleges for search", len(collegesResponse))
	return collegesResponse
}
