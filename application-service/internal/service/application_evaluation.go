package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/compendium-tech/compendium/application-service/internal/context"
	"github.com/compendium-tech/compendium/application-service/internal/domain"
	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/compendium-tech/compendium/application-service/internal/repository"
	llmdomain "github.com/compendium-tech/compendium/llm-common/pkg/domain"
	"github.com/compendium-tech/compendium/llm-common/pkg/service"
)

type ApplicationEvaluationService interface {
	EvaluateCurrentApplication(ctx context.Context) (*domain.ApplicationEvaluationResponse, error)
}

type applicationEvaluationService struct {
	applicationRepository repository.ApplicationRepository
	llmService            service.LLMService
}

func NewApplicationEvaluateService(
	applicationRepository repository.ApplicationRepository,
	llmService service.LLMService) ApplicationEvaluationService {
	return &applicationEvaluationService{
		applicationRepository: applicationRepository,
		llmService:            llmService,
	}
}

func (s *applicationEvaluationService) EvaluateCurrentApplication(ctx context.Context) (*domain.ApplicationEvaluationResponse, error) {
	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return nil, err
	}

	activities, err := s.applicationRepository.GetActivities(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	honors, err := s.applicationRepository.GetHonors(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	essays, err := s.applicationRepository.GetEssays(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	supplementalEssays, err := s.applicationRepository.GetSupplementalEssays(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	activitiesEvaluation, err := s.evaluateActivities(ctx, activities)
	if err != nil {
		return nil, err
	}

	honorsEvaluation, err := s.evaluateHonors(ctx, honors)
	if err != nil {
		return nil, err
	}

	essaysEvaluation, err := s.evaluateEssays(ctx, essays)
	if err != nil {
		return nil, err
	}

	supplementalEssaysEvaluation, err := s.evaluateSupplementalEssays(ctx, supplementalEssays)
	if err != nil {
		return nil, err
	}

	return &domain.ApplicationEvaluationResponse{
		ActivitiesEvaluationResponse:         *activitiesEvaluation,
		HonorsEvaluationResponse:             *honorsEvaluation,
		EssaysEvaluationResponse:             *essaysEvaluation,
		SupplementalEssaysEvaluationResponse: *supplementalEssaysEvaluation,
	}, nil
}

func (s *applicationEvaluationService) evaluateActivities(ctx context.Context, activities []model.Activity) (*domain.ActivitiesEvaluationResponse, error) {
	prompt := activitiesEvaluationPromptBase

	for idx, activity := range activities {
		prompt += fmt.Sprintf("%d. %s - %s\n", idx+1, activity.Role, activity.Name)

		if activity.Description != nil {
			prompt += "Description: " + *activity.Description
		}

		prompt += fmt.Sprintf("Hours per week: %d\n", activity.HoursPerWeek)
		prompt += fmt.Sprintf("Weeks per year: %d\n", activity.WeeksPerYear)
		prompt += fmt.Sprintf("Category: %s\n", activity.Category)

		gradeStrings := make([]string, len(activity.Grades))
		for _, grade := range activity.Grades {
			gradeStrings = append(gradeStrings, string(grade))
		}

		prompt += fmt.Sprintf("Grade levels: %s\n", strings.Join(gradeStrings, ", "))
	}

	llmResponse, err := s.llmService.GenerateResponse(ctx, []llmdomain.Message{
		{
			Role: llmdomain.RoleSystem,
			Text: prompt,
		},
	}, nil, &activitiesEvaluationSchema)
	if err != nil {
		return nil, err
	}

	var response domain.ActivitiesEvaluationResponse
	err = json.Unmarshal([]byte(llmResponse.Text), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *applicationEvaluationService) evaluateHonors(ctx context.Context, honors []model.Honor) (*domain.HonorsEvaluationResponse, error) {
	prompt := honorsEvaluationPromptBase

	for idx, honor := range honors {
		prompt += fmt.Sprintf("%d. %s - %s\n", idx+1, honor.Title)

		if honor.Description != nil {
			prompt += "Description: " + *honor.Description
		}

		prompt += fmt.Sprintf("Level: %s\n", honor.Level)
		prompt += fmt.Sprintf("Grade: %s\n", honor.Grade)
	}

	llmResponse, err := s.llmService.GenerateResponse(ctx, []llmdomain.Message{
		{
			Role: llmdomain.RoleSystem,
			Text: prompt,
		},
	}, nil, &honorsEvaluationSchema)
	if err != nil {
		return nil, err
	}

	var response domain.HonorsEvaluationResponse
	err = json.Unmarshal([]byte(llmResponse.Text), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *applicationEvaluationService) evaluateEssays(ctx context.Context, essays []model.Essay) (*domain.EssaysEvaluationResponse, error) {
	prompt := essaysEvaluationPromptBase

	for idx, essay := range essays {
		prompt += fmt.Sprintf("%d. %s (type: %s)", idx+1, essay.Type)
		prompt += essay.Content + "\n\n\n"
	}

	structuredOutputSchema := generateEssaysEvaluationSchema(len(essays))

	llmResponse, err := s.llmService.GenerateResponse(ctx, []llmdomain.Message{
		{
			Role: llmdomain.RoleSystem,
			Text: prompt,
		},
	}, nil, &structuredOutputSchema)
	if err != nil {
		return nil, err
	}

	var response domain.EssaysEvaluationResponse
	err = json.Unmarshal([]byte(llmResponse.Text), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *applicationEvaluationService) evaluateSupplementalEssays(ctx context.Context, essays []model.SupplementalEssay) (*domain.SupplementalEssaysEvaluationResponse, error) {
	prompt := supplementalEssaysEvaluationPromptBase

	for idx, essay := range essays {
		prompt += fmt.Sprintf("%d. %s\n", idx+1, essay.Title)
		prompt += essay.Content + "\n\n\n"
	}

	structuredOutputSchema := generateSupplementalEssaysEvaluationSchema(len(essays))

	llmResponse, err := s.llmService.GenerateResponse(ctx, []llmdomain.Message{
		{
			Role: llmdomain.RoleSystem,
			Text: prompt,
		},
	}, nil, &structuredOutputSchema)
	if err != nil {
		return nil, err
	}

	var response domain.SupplementalEssaysEvaluationResponse
	err = json.Unmarshal([]byte(llmResponse.Text), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
