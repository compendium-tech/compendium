package service

import (
	"context"
	"time"

	"github.com/compendium-tech/compendium/application-assistant-service/internal/domain"
	appErr "github.com/compendium-tech/compendium/application-assistant-service/internal/error"
	"github.com/compendium-tech/compendium/application-assistant-service/internal/model"
	"github.com/compendium-tech/compendium/application-assistant-service/internal/repository"
	"github.com/compendium-tech/compendium/common/pkg/auth"
	log "github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/google/uuid"
)

type ApplicationService interface {
	GetApplications(ctx context.Context) ([]domain.ApplicationResponse, error)
	CreateApplication(ctx context.Context, request domain.CreateApplicationRequest) (*domain.ApplicationResponse, error)
	RemoveApplication(ctx context.Context, id uuid.UUID) error

	AddActivity(ctx context.Context, applicationID uuid.UUID, activity domain.UpdateActivityRequest) (*domain.ActivityResponse, error)
	RemoveActivity(ctx context.Context, applicationID uuid.UUID) error
	UpdateActivity(ctx context.Context, applicationID uuid.UUID, request domain.UpdateActivityRequest) (*domain.ActivityResponse, error)
	GetActivities(ctx context.Context, applicationID uuid.UUID) ([]domain.ActivityResponse, error)

	AddHonor(ctx context.Context, applicationID uuid.UUID, honor domain.UpdateHonorRequest) (*domain.HonorResponse, error)
	RemoveHonor(ctx context.Context, applicationID uuid.UUID) error
	UpdateHonor(ctx context.Context, applicationID uuid.UUID, request domain.UpdateHonorRequest) (*domain.HonorResponse, error)
	GetHonors(ctx context.Context, applicationID uuid.UUID) ([]domain.HonorResponse, error)
}

type applicationService struct {
	applicationRepository repository.ApplicationRepository
}

func NewApplicationService(applicationRepository repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		applicationRepository: applicationRepository,
	}
}

func (s *applicationService) GetApplications(ctx context.Context) ([]domain.ApplicationResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String())
	logger.Info("Fetching applications for authenticated user")

	applications, err := s.applicationRepository.FindApplicationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	logger.Info("Applications fetched successfully")

	responses := make([]domain.ApplicationResponse, len(applications))
	for i, app := range applications {
		responses[i] = domain.ApplicationResponse{
			ID:        app.ID,
			Name:      app.Name,
			CreatedAt: app.CreatedAt,
		}
	}
	return responses, nil
}

func (s *applicationService) CreateApplication(ctx context.Context, request domain.CreateApplicationRequest) (*domain.ApplicationResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("name", request.Name)
	logger.Info("Creating new application")

	app := model.Application{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      request.Name,
		CreatedAt: time.Now(),
	}

	if err := s.applicationRepository.CreateApplication(ctx, app); err != nil {
		return nil, err
	}

	logger.Info("Application created successfully")

	return &domain.ApplicationResponse{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
	}, nil
}

func (s *applicationService) RemoveApplication(ctx context.Context, id uuid.UUID) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", id.String())
	logger.Info("Deleting application")

	app, err := s.applicationRepository.GetApplication(ctx, id)
	if err != nil {
		return err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	if err := s.applicationRepository.RemoveApplication(ctx, id); err != nil {
		return err
	}

	logger.Info("Application deleted successfully")
	return nil
}

func (s *applicationService) AddActivity(ctx context.Context, applicationID uuid.UUID, request domain.UpdateActivityRequest) (*domain.ActivityResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Adding activity to application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return nil, appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	activity := model.Activity{
		ApplicationID: applicationID,
		Name:          request.Name,
		Role:          request.Role,
		Description:   request.Description,
		HoursPerWeek:  request.HoursPerWeek,
		WeeksPerYear:  request.WeeksPerYear,
		Category:      request.Category,
		Grades:        request.Grades,
	}

	if err := s.applicationRepository.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}

	logger.Info("Activity added successfully")

	return &domain.ActivityResponse{
		Name:         activity.Name,
		Role:         activity.Role,
		Description:  activity.Description,
		HoursPerWeek: activity.HoursPerWeek,
		WeeksPerYear: activity.WeeksPerYear,
		Category:     activity.Category,
		Grades:       activity.Grades,
	}, nil
}

func (s *applicationService) RemoveActivity(ctx context.Context, applicationID uuid.UUID) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Removing activity from application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	if err := s.applicationRepository.DeleteActivity(ctx, applicationID); err != nil {
		return err
	}

	logger.Info("Activity removed successfully")
	return nil
}

func (s *applicationService) UpdateActivity(ctx context.Context, applicationID uuid.UUID, request domain.UpdateActivityRequest) (*domain.ActivityResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Updating activity for application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return nil, appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	activity := model.Activity{
		ApplicationID: applicationID,
		Name:          request.Name,
		Role:          request.Role,
		Description:   request.Description,
		HoursPerWeek:  request.HoursPerWeek,
		WeeksPerYear:  request.WeeksPerYear,
		Category:      request.Category,
		Grades:        request.Grades,
	}

	updatedActivity, err := s.applicationRepository.UpdateActivity(ctx, activity)
	if err != nil {
		return nil, err
	}

	logger.Info("Activity updated successfully")

	return &domain.ActivityResponse{
		Name:         updatedActivity.Name,
		Role:         updatedActivity.Role,
		Description:  updatedActivity.Description,
		HoursPerWeek: updatedActivity.HoursPerWeek,
		WeeksPerYear: updatedActivity.WeeksPerYear,
		Category:     updatedActivity.Category,
		Grades:       updatedActivity.Grades,
	}, nil
}

func (s *applicationService) GetActivities(ctx context.Context, applicationID uuid.UUID) ([]domain.ActivityResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Fetching activities for application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found or unauthorized access")
		return nil, appErr.New(appErr.ApplicationNotFoundError, "application not found or unauthorized")
	}

	activities, err := s.applicationRepository.GetActivities(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	logger.Info("Activities fetched successfully")

	responses := make([]domain.ActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = domain.ActivityResponse{
			Name:         activity.Name,
			Role:         activity.Role,
			Description:  activity.Description,
			HoursPerWeek: activity.HoursPerWeek,
			WeeksPerYear: activity.WeeksPerYear,
			Category:     activity.Category,
			Grades:       activity.Grades,
		}
	}

	return responses, nil
}

func (s *applicationService) AddHonor(ctx context.Context, activityID uuid.UUID, request domain.UpdateHonorRequest) (*domain.HonorResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", activityID.String())
	logger.Info("Adding honor to application")

	app, err := s.applicationRepository.GetApplication(ctx, activityID)
	if err != nil {
		return nil, err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return nil, appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	honor := model.Honor{
		Title:       request.Title,
		Description: request.Description,
		Level:       request.Level,
		Grade:       request.Grade,
	}

	if err := s.applicationRepository.CreateHonor(ctx, honor); err != nil {
		return nil, err
	}

	logger.Info("Honor added successfully")

	return &domain.HonorResponse{
		Title:       honor.Title,
		Description: honor.Description,
		Level:       honor.Level,
		Grade:       honor.Grade,
	}, nil
}

func (s *applicationService) RemoveHonor(ctx context.Context, applicationID uuid.UUID) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Removing honor from application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	if err := s.applicationRepository.DeleteHonor(ctx, applicationID); err != nil {
		return err
	}

	logger.Info("Honor removed successfully")
	return nil
}

func (s *applicationService) UpdateHonor(ctx context.Context, applicationID uuid.UUID, request domain.UpdateHonorRequest) (*domain.HonorResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Updating honor for application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found")
		return nil, appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	honor := model.Honor{
		ApplicationID: applicationID,
		Title:         request.Title,
		Description:   request.Description,
		Level:         request.Level,
		Grade:         request.Grade,
	}

	updatedHonor, err := s.applicationRepository.UpdateHonor(ctx, honor)
	if err != nil {
		return nil, err
	}

	logger.Info("Honor updated successfully")

	return &domain.HonorResponse{
		Title:       updatedHonor.Title,
		Description: updatedHonor.Description,
		Level:       updatedHonor.Level,
		Grade:       updatedHonor.Grade,
	}, nil
}

func (s *applicationService) GetHonors(ctx context.Context, applicationID uuid.UUID) ([]domain.HonorResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.L(ctx).WithField("userId", userID.String()).WithField("applicationId", applicationID.String())
	logger.Info("Fetching honors for application")

	app, err := s.applicationRepository.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app == nil || app.UserID != userID {
		logger.Warn("Application not found or unauthorized access")
		return nil, appErr.New(appErr.ApplicationNotFoundError, "application not found or unauthorized")
	}

	honors, err := s.applicationRepository.GetHonors(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	logger.Info("Honors fetched successfully")

	responses := make([]domain.HonorResponse, len(honors))
	for i, honor := range honors {
		responses[i] = domain.HonorResponse{
			Title:       honor.Title,
			Description: honor.Description,
			Level:       honor.Level,
			Grade:       honor.Grade,
		}
	}

	return responses, nil
}
