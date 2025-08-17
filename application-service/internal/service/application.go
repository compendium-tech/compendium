package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"

	localcontext "github.com/compendium-tech/compendium/application-service/internal/context"
	"github.com/compendium-tech/compendium/application-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/application-service/internal/error"
	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/compendium-tech/compendium/application-service/internal/repository"
)

type ApplicationService interface {
	GetCurrentApplicationModel(ctx context.Context, id uuid.UUID) model.Application

	GetApplications(ctx context.Context) []domain.ApplicationResponse
	CreateApplication(ctx context.Context, request domain.CreateApplicationRequest) domain.ApplicationResponse

	UpdateCurrentApplicationName(ctx context.Context, name string)
	RemoveCurrentApplication(ctx context.Context)

	GetActivities(ctx context.Context) []domain.ActivityResponse
	PutActivities(ctx context.Context, activities []domain.UpdateActivityRequest)
	GetHonors(ctx context.Context) []domain.HonorResponse
	PutHonors(ctx context.Context, honors []domain.UpdateHonorRequest)

	GetEssays(ctx context.Context) []domain.EssayResponse
	PutEssays(ctx context.Context, essays []domain.UpdateEssayRequest)
	GetSupplementalEssays(ctx context.Context) []domain.SupplementalEssayResponse
	PutSupplementalEssays(ctx context.Context, supplementalEssays []domain.UpdateSupplementalEssayRequest)
}

type applicationService struct {
	applicationRepository repository.ApplicationRepository
}

func NewApplicationService(applicationRepository repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		applicationRepository: applicationRepository,
	}
}

func (a *applicationService) GetApplications(ctx context.Context) []domain.ApplicationResponse {
	userID := auth.GetUserID(ctx)
	log.L(ctx).Info("Fetching applications")

	applications := a.applicationRepository.FindApplicationsByUserID(ctx, userID)
	applicationResponses := make([]domain.ApplicationResponse, len(applications))
	for i, application := range applications {
		applicationResponses[i] = domain.ApplicationResponse{
			ID:        application.ID,
			Name:      application.Name,
			CreatedAt: application.CreatedAt,
		}
	}

	log.L(ctx).Infof("Found %d applications", len(applicationResponses))
	return applicationResponses
}

func (a *applicationService) GetCurrentApplicationModel(ctx context.Context, id uuid.UUID) model.Application {
	userID := auth.GetUserID(ctx)
	logger := log.L(ctx).WithField("applicationId", id.String())
	logger.Info("Getting current application model")

	application := a.applicationRepository.GetApplication(ctx, id)
	if application == nil || application.UserID != userID {
		logger.Warn("Application not found")
		myerror.New(myerror.ApplicationNotFoundError).Throw()
	}

	logger.Info("Application model fetched successfully")
	return *application
}

func (a *applicationService) CreateApplication(ctx context.Context, request domain.CreateApplicationRequest) domain.ApplicationResponse {
	userID := auth.GetUserID(ctx)
	logger := log.L(ctx).WithField("applicationName", request.Name)
	logger.Info("Creating application")

	id := uuid.New()
	createdAt := time.Now().UTC()
	a.applicationRepository.CreateApplication(ctx, model.Application{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      request.Name,
		CreatedAt: createdAt,
	})

	logger.Info("Application created successfully")
	return domain.ApplicationResponse{
		ID:        id,
		Name:      request.Name,
		CreatedAt: createdAt,
	}
}

func (a *applicationService) UpdateCurrentApplicationName(ctx context.Context, name string) {
	logger := log.L(ctx).WithField("newName", name)
	logger.Info("Updating current application name")

	a.applicationRepository.UpdateApplicationName(
		ctx, localcontext.GetApplication(ctx).ID, name)
	logger.Info("Application name updated successfully")
}

func (a *applicationService) RemoveCurrentApplication(ctx context.Context) {
	log.L(ctx).Info("Removing current application")

	a.applicationRepository.RemoveApplication(
		ctx, localcontext.GetApplication(ctx).ID)
	log.L(ctx).Info("Application removed successfully")
}

func (a *applicationService) GetActivities(ctx context.Context) []domain.ActivityResponse {
	log.L(ctx).Info("Getting activities")

	activities := a.applicationRepository.
		GetActivities(ctx, localcontext.GetApplication(ctx).ID)
	activitiesResponse := make([]domain.ActivityResponse, len(activities))

	for i, activity := range activities {
		activitiesResponse[i] = domain.ActivityResponse{
			ID:           activity.ID,
			Name:         activity.Name,
			Role:         activity.Role,
			Description:  activity.Description,
			HoursPerWeek: activity.HoursPerWeek,
			WeeksPerYear: activity.WeeksPerYear,
			Category:     activity.Category,
			Grades:       activity.Grades,
		}
	}

	log.L(ctx).Infof("Found %d activities", len(activitiesResponse))
	return activitiesResponse
}

func (a *applicationService) PutActivities(ctx context.Context, updateActivitiesRequest []domain.UpdateActivityRequest) {
	logger := log.L(ctx).WithField("activityCount", len(updateActivitiesRequest))
	logger.Info("Putting activities")

	application := localcontext.GetApplication(ctx)
	activities := make([]model.Activity, len(updateActivitiesRequest))

	for i, updateActivityRequest := range updateActivitiesRequest {
		activities[i] = model.Activity{
			ID:           uuid.New(),
			Name:         updateActivityRequest.Name,
			Role:         updateActivityRequest.Role,
			Description:  updateActivityRequest.Description,
			HoursPerWeek: updateActivityRequest.HoursPerWeek,
			WeeksPerYear: updateActivityRequest.WeeksPerYear,
			Category:     updateActivityRequest.Category,
			Grades:       updateActivityRequest.Grades,
		}
	}

	a.applicationRepository.PutActivities(ctx, application.ID, activities)
	logger.Info("Activities put successfully")
}

func (a *applicationService) GetHonors(ctx context.Context) []domain.HonorResponse {
	log.L(ctx).Info("Getting honors")

	honors := a.applicationRepository.
		GetHonors(ctx, localcontext.GetApplication(ctx).ID)
	honorsResponse := make([]domain.HonorResponse, len(honors))

	for i, honor := range honors {
		honorsResponse[i] = domain.HonorResponse{
			Title:       honor.Title,
			Description: honor.Description,
			Level:       honor.Level,
			Grade:       honor.Grade,
		}
	}

	log.L(ctx).Infof("Found %d honors", len(honorsResponse))
	return honorsResponse
}

func (a *applicationService) PutHonors(ctx context.Context, updateHonorsRequest []domain.UpdateHonorRequest) {
	logger := log.L(ctx).WithField("honorCount", len(updateHonorsRequest))
	logger.Info("Putting honors")

	application := localcontext.GetApplication(ctx)
	honors := make([]model.Honor, len(updateHonorsRequest))

	for i, updateHonorRequest := range updateHonorsRequest {
		honors[i] = model.Honor{
			ID:          uuid.New(),
			Title:       updateHonorRequest.Title,
			Description: updateHonorRequest.Description,
			Level:       updateHonorRequest.Level,
			Grade:       updateHonorRequest.Grade,
		}
	}

	a.applicationRepository.PutHonors(ctx, application.ID, honors)
	logger.Info("Honors put successfully")
}

func (a *applicationService) GetEssays(ctx context.Context) []domain.EssayResponse {
	log.L(ctx).Info("Getting essays")

	essays := a.applicationRepository.
		GetEssays(ctx, localcontext.GetApplication(ctx).ID)
	essaysResponse := make([]domain.EssayResponse, len(essays))

	for i, essay := range essays {
		essaysResponse[i] = domain.EssayResponse{
			ID:      essay.ID,
			Kind:    essay.Type,
			Content: essay.Content,
		}
	}

	log.L(ctx).Infof("Found %d essays", len(essaysResponse))
	return essaysResponse
}

func (a *applicationService) PutEssays(ctx context.Context, updateEssaysRequest []domain.UpdateEssayRequest) {
	logger := log.L(ctx).WithField("essayCount", len(updateEssaysRequest))
	logger.Info("Putting essays")

	application := localcontext.GetApplication(ctx)
	essays := make([]model.Essay, len(updateEssaysRequest))

	for i, updateEssayRequest := range updateEssaysRequest {
		essays[i] = model.Essay{
			ID:      uuid.New(),
			Type:    updateEssayRequest.Kind,
			Content: updateEssayRequest.Content,
		}
	}

	a.applicationRepository.PutEssays(ctx, application.ID, essays)
	logger.Info("Essays put successfully")
}

func (a *applicationService) GetSupplementalEssays(ctx context.Context) []domain.SupplementalEssayResponse {
	log.L(ctx).Info("Getting supplemental essays")

	supplementalEssays := a.applicationRepository.
		GetSupplementalEssays(ctx, localcontext.GetApplication(ctx).ID)
	supplementalEssaysResponse := make([]domain.SupplementalEssayResponse, len(supplementalEssays))

	for i, supplementalEssay := range supplementalEssays {
		supplementalEssaysResponse[i] = domain.SupplementalEssayResponse{
			ID:      supplementalEssay.ID,
			Title:   supplementalEssay.Prompt,
			Content: supplementalEssay.Content,
		}
	}

	log.L(ctx).Infof("Found %d supplemental essays", len(supplementalEssaysResponse))
	return supplementalEssaysResponse
}

func (a *applicationService) PutSupplementalEssays(ctx context.Context, updateSupplementalEssaysRequest []domain.UpdateSupplementalEssayRequest) {
	logger := log.L(ctx).WithField("supplementalEssayCount", len(updateSupplementalEssaysRequest))
	logger.Info("Putting supplemental essays")

	application := localcontext.GetApplication(ctx)
	supplementalEssays := make([]model.SupplementalEssay, len(updateSupplementalEssaysRequest))

	for _, updateSupplementalEssayRequest := range updateSupplementalEssaysRequest {
		supplementalEssays = append(supplementalEssays, model.SupplementalEssay{
			ID:      uuid.New(),
			Prompt:  updateSupplementalEssayRequest.Title,
			Content: updateSupplementalEssayRequest.Content,
		})
	}

	a.applicationRepository.PutSupplementalEssays(ctx, application.ID, supplementalEssays)
	logger.Info("Supplemental essays put successfully")
}
