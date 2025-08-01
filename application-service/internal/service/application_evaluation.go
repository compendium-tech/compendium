package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	localcontext "github.com/compendium-tech/compendium/application-service/internal/context"
	"github.com/compendium-tech/compendium/application-service/internal/domain"
	"github.com/compendium-tech/compendium/application-service/internal/interop"
	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/compendium-tech/compendium/application-service/internal/repository"
)

type ApplicationEvaluationService interface {
	EvaluateCurrentApplication(ctx context.Context) (*domain.ApplicationEvaluationResponse, error)
}

type applicationEvaluationService struct {
	applicationRepository repository.ApplicationRepository
	llmService            interop.LLMService
}

func NewApplicationEvaluateService(
	applicationRepository repository.ApplicationRepository,
	llmService interop.LLMService) ApplicationEvaluationService {
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

	return s.evaluateApplication(ctx, activities, honors, essays, supplementalEssays)
}

func (s *applicationEvaluationService) evaluateApplication(
	ctx context.Context, activities []model.Activity,
	honors []model.Honor, essays []model.Essay,
	supplementalEssays []model.SupplementalEssay) (*domain.ApplicationEvaluationResponse, error) {
	prompt := applicationEvaluationPromptBase
	structuredOutputSchema := generateApplicationEvaluationSchema(len(essays), len(supplementalEssays))

	prompt += "# Application to evaluate\n\n## Extracurricular activities"
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

	prompt += "## Honors\n"
	for idx, honor := range honors {
		prompt += fmt.Sprintf("%d. %s\n", idx+1, honor.Title)

		if honor.Description != nil {
			prompt += "Description: " + *honor.Description
		}

		prompt += fmt.Sprintf("Level: %s\n", honor.Level)
		prompt += fmt.Sprintf("Grade: %s\n", honor.Grade)
	}

	prompt += "## Essays\n"
	for idx, essay := range essays {
		prompt += fmt.Sprintf("%d. Type: %s", idx+1, essay.Type)
		prompt += essay.Content + "\n\n\n"
	}

	prompt += "## Supplemental essays\n"
	for idx, essay := range supplementalEssays {
		prompt += fmt.Sprintf("%d. Prompt: %s\n", idx+1, essay.Prompt)
		prompt += essay.Content + "\n\n\n"
	}

	llmResponse, err := s.llmService.GenerateResponse(ctx, []domain.LLMMessage{
		{
			Role: domain.RoleSystem,
			Text: prompt,
		},
	}, nil, &structuredOutputSchema)
	if err != nil {
		return nil, err
	}

	var response domain.ApplicationEvaluationResponse
	err = json.Unmarshal([]byte(llmResponse.Text), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
