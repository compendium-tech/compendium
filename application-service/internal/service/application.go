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
	GetCurrentApplicationModel(ctx context.Context, id uuid.UUID) (*model.Application, error)

	GetApplications(ctx context.Context) ([]domain.ApplicationResponse, error)
	CreateApplication(ctx context.Context, request domain.CreateApplicationRequest) (*domain.ApplicationResponse, error)

	UpdateCurrentApplicationName(ctx context.Context, name string) error
	RemoveCurrentApplication(ctx context.Context) error

	GetActivities(ctx context.Context) ([]domain.ActivityResponse, error)
	PutActivities(ctx context.Context, activities []domain.UpdateActivityRequest) error
	GetHonors(ctx context.Context) ([]domain.HonorResponse, error)
	PutHonors(ctx context.Context, honors []domain.UpdateHonorRequest) error

	GetEssays(ctx context.Context) ([]domain.EssayResponse, error)
	PutEssays(ctx context.Context, essays []domain.UpdateEssayRequest) error
	GetSupplementalEssays(ctx context.Context) ([]domain.SupplementalEssayResponse, error)
	PutSupplementalEssays(ctx context.Context, supplementalEssays []domain.UpdateSupplementalEssayRequest) error
}

type applicationService struct {
	applicationRepository repository.ApplicationRepository
}

func NewApplicationService(applicationRepository repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		applicationRepository: applicationRepository,
	}
}

func (a *applicationService) GetApplications(ctx context.Context) ([]domain.ApplicationResponse, error) {
	log.L(ctx).Info("Getting applications")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	applications, err := a.applicationRepository.FindApplicationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var applicationResponses []domain.ApplicationResponse
	for _, application := range applications {
		applicationResponses = append(applicationResponses, domain.ApplicationResponse{
			ID:        application.ID,
			Name:      application.Name,
			CreatedAt: application.CreatedAt,
		})
	}

	log.L(ctx).Infof("Found %d applications", len(applicationResponses))
	return applicationResponses, nil
}

func (a *applicationService) GetCurrentApplicationModel(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	logger := log.L(ctx).WithField("applicationId", id.String())
	logger.Info("Getting current application model")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	application, err := a.applicationRepository.GetApplication(ctx, id)
	if err != nil {
		return nil, err
	}

	if application == nil || application.UserID != userID {
		logger.Warn("Application not found")
		return nil, myerror.New(myerror.ApplicationNotFoundError)
	}

	logger.Info("Application model fetched successfully")
	return application, nil
}

func (a *applicationService) CreateApplication(ctx context.Context, request domain.CreateApplicationRequest) (*domain.ApplicationResponse, error) {
	logger := log.L(ctx).WithField("applicationName", request.Name)
	logger.Info("Creating application")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	application := model.Application{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      request.Name,
		CreatedAt: time.Now().UTC(),
	}

	err = a.applicationRepository.CreateApplication(ctx, application)
	if err != nil {
		return nil, err
	}

	logger.Info("Application created successfully")
	return &domain.ApplicationResponse{
		ID:        application.ID,
		Name:      application.Name,
		CreatedAt: application.CreatedAt,
	}, nil
}

func (a *applicationService) UpdateCurrentApplicationName(ctx context.Context, name string) error {
	logger := log.L(ctx).WithField("newName", name)
	logger.Info("Updating current application name")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return err
	}

	err = a.applicationRepository.UpdateApplicationName(ctx, application.ID, name)
	if err != nil {
		return err
	}

	logger.Info("Application name updated successfully")
	return nil
}

func (a *applicationService) RemoveCurrentApplication(ctx context.Context) error {
	log.L(ctx).Info("Removing current application")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return err
	}

	err = a.applicationRepository.RemoveApplication(ctx, application.ID)
	if err != nil {
		return err
	}

	log.L(ctx).Info("Application removed successfully")
	return nil
}

func (a *applicationService) GetActivities(ctx context.Context) ([]domain.ActivityResponse, error) {
	log.L(ctx).Info("Getting activities")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return nil, err
	}

	activities, err := a.applicationRepository.GetActivities(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	var activityResponses []domain.ActivityResponse
	for _, activity := range activities {
		activityResponses = append(activityResponses, domain.ActivityResponse{
			ID:           activity.ID,
			Name:         activity.Name,
			Role:         activity.Role,
			Description:  activity.Description,
			HoursPerWeek: activity.HoursPerWeek,
			WeeksPerYear: activity.WeeksPerYear,
			Category:     activity.Category,
			Grades:       activity.Grades,
		})
	}

	log.L(ctx).Infof("Found %d activities", len(activityResponses))
	return activityResponses, nil
}

func (a *applicationService) PutActivities(ctx context.Context, activityRequests []domain.UpdateActivityRequest) error {
	logger := log.L(ctx).WithField("activityCount", len(activityRequests))
	logger.Info("Putting activities")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return err
	}

	var activities []model.Activity
	for _, activityRequest := range activityRequests {
		activities = append(activities, model.Activity{
			ID:           uuid.New(),
			Name:         activityRequest.Name,
			Role:         activityRequest.Role,
			Description:  activityRequest.Description,
			HoursPerWeek: activityRequest.HoursPerWeek,
			WeeksPerYear: activityRequest.WeeksPerYear,
			Category:     activityRequest.Category,
			Grades:       activityRequest.Grades,
		})
	}

	err = a.applicationRepository.PutActivities(ctx, application.ID, activities)
	if err != nil {
		return err
	}

	logger.Info("Activities put successfully")
	return nil
}

func (a *applicationService) GetHonors(ctx context.Context) ([]domain.HonorResponse, error) {
	log.L(ctx).Info("Getting honors")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return nil, err
	}

	honors, err := a.applicationRepository.GetHonors(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	var honorResponses []domain.HonorResponse
	for _, honor := range honors {
		honorResponses = append(honorResponses, domain.HonorResponse{
			Title:       honor.Title,
			Description: honor.Description,
			Level:       honor.Level,
			Grade:       honor.Grade,
		})
	}

	log.L(ctx).Infof("Found %d honors", len(honorResponses))
	return honorResponses, nil
}

func (a *applicationService) PutHonors(ctx context.Context, honorRequests []domain.UpdateHonorRequest) error {
	logger := log.L(ctx).WithField("honorCount", len(honorRequests))
	logger.Info("Putting honors")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return err
	}

	var honors []model.Honor
	for _, honorRequest := range honorRequests {
		honors = append(honors, model.Honor{
			ID:          uuid.New(),
			Title:       honorRequest.Title,
			Description: honorRequest.Description,
			Level:       honorRequest.Level,
			Grade:       honorRequest.Grade,
		})
	}

	err = a.applicationRepository.PutHonors(ctx, application.ID, honors)
	if err != nil {
		return err
	}

	logger.Info("Honors put successfully")
	return nil
}

func (a *applicationService) GetEssays(ctx context.Context) ([]domain.EssayResponse, error) {
	log.L(ctx).Info("Getting essays")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return nil, err
	}

	essays, err := a.applicationRepository.GetEssays(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	var essayResponses []domain.EssayResponse
	for _, essay := range essays {
		essayResponses = append(essayResponses, domain.EssayResponse{
			ID:      essay.ID,
			Kind:    essay.Type,
			Content: essay.Content,
		})
	}

	log.L(ctx).Infof("Found %d essays", len(essayResponses))
	return essayResponses, nil
}

func (a *applicationService) PutEssays(ctx context.Context, essayRequests []domain.UpdateEssayRequest) error {
	logger := log.L(ctx).WithField("essayCount", len(essayRequests))
	logger.Info("Putting essays")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return err
	}

	var essays []model.Essay
	for _, essayRequest := range essayRequests {
		essays = append(essays, model.Essay{
			ID:      uuid.New(),
			Type:    essayRequest.Kind,
			Content: essayRequest.Content,
		})
	}

	err = a.applicationRepository.PutEssays(ctx, application.ID, essays)
	if err != nil {
		return err
	}

	logger.Info("Essays put successfully")
	return nil
}

func (a *applicationService) GetSupplementalEssays(ctx context.Context) ([]domain.SupplementalEssayResponse, error) {
	log.L(ctx).Info("Getting supplemental essays")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return nil, err
	}

	supplementalEssays, err := a.applicationRepository.GetSupplementalEssays(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	var supplementalEssayResponses []domain.SupplementalEssayResponse
	for _, supplementalEssay := range supplementalEssays {
		supplementalEssayResponses = append(supplementalEssayResponses, domain.SupplementalEssayResponse{
			ID:      supplementalEssay.ID,
			Title:   supplementalEssay.Prompt,
			Content: supplementalEssay.Content,
		})
	}

	log.L(ctx).Infof("Found %d supplemental essays", len(supplementalEssayResponses))
	return supplementalEssayResponses, nil
}

func (a *applicationService) PutSupplementalEssays(ctx context.Context, supplementalEssayRequests []domain.UpdateSupplementalEssayRequest) error {
	logger := log.L(ctx).WithField("supplementalEssayCount", len(supplementalEssayRequests))
	logger.Info("Putting supplemental essays")

	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return err
	}

	var supplementalEssays []model.SupplementalEssay
	for _, supplementalEssayRequest := range supplementalEssayRequests {
		supplementalEssays = append(supplementalEssays, model.SupplementalEssay{
			ID:      uuid.New(),
			Prompt:  supplementalEssayRequest.Title,
			Content: supplementalEssayRequest.Content,
		})
	}

	err = a.applicationRepository.PutSupplementalEssays(ctx, application.ID, supplementalEssays)
	if err != nil {
		return err
	}

	logger.Info("Supplemental essays put successfully")
	return nil
}
