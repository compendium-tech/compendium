package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compendium-tech/compendium/application-service/internal/context"
	"github.com/compendium-tech/compendium/application-service/internal/domain"
	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/compendium-tech/compendium/application-service/internal/repository"
	llmdomain "github.com/compendium-tech/compendium/llm-common/pkg/domain"
	"github.com/compendium-tech/compendium/llm-common/pkg/service"
	"strings"
)

type ApplicationAssessmentService interface {
	EvaluateCurrentApplication(ctx context.Context) (*domain.ApplicationEvaluationResponse, error)
}

type applicationAssessmentService struct {
	applicationRepository repository.ApplicationRepository
	llmService            service.LLMService
}

func NewApplicationAssessmentService(
	applicationRepository repository.ApplicationRepository,
	llmService service.LLMService) ApplicationAssessmentService {
	return &applicationAssessmentService{
		applicationRepository: applicationRepository,
		llmService:            llmService,
	}
}

func (s *applicationAssessmentService) EvaluateCurrentApplication(ctx context.Context) (*domain.ApplicationEvaluationResponse, error) {
	application, err := localcontext.GetApplication(ctx)
	if err != nil {
		return nil, err
	}

	activities, err := s.applicationRepository.GetActivities(ctx, application.ID)
	if err != nil {
		return nil, err
	}

	activitiesEvaluation, err := s.EvaluateActivities(ctx, activities)
	if err != nil {
		return nil, err
	}

	return &domain.ApplicationEvaluationResponse{
		ActivitiesEvaluationResponse: *activitiesEvaluation,
	}, nil
}

func (s *applicationAssessmentService) EvaluateActivities(
	ctx context.Context, activities []model.Activity) (*domain.ActivitiesEvaluationResponse, error) {
	prompt := applicationEvaluationPromptBase

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
	}, &applicationEvaluationSchema)
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

var applicationEvaluationSchema = llmdomain.StructuredOutputSchema{
	Type: llmdomain.TypeObject,
	Properties: map[string]llmdomain.StructuredOutputSchema{
		"suggestions": {
			Type: llmdomain.TypeString,
			Description: `A list of actionable recommendations to improve the activities section, such 
as reordering activities, adding specific achievements, or increasing involvement in relevant activities.`,
		},
		"strengths": {
			Type: llmdomain.TypeString,
			Description: `A list of key strengths in the activities section, highlighting deep involvement, 
leadership roles, impactful contributions, or alignment with the student’s goals.`,
		},
		"weaknesses": {
			Type: llmdomain.TypeString,
			Description: `A list of weaknesses in the activities section, such as superficial involvement, lack 
of leadership, unclear descriptions, or misalignment with the student’s goals.`,
		},
		"summary": {
			Type: llmdomain.TypeString,
			Description: `A concise summary of the overall quality of the activities section, addressing depth, 
impact, relevance, and presentation, with an evaluation of how well it reflects the student’s strengths and goals.`,
		},
	},
	Required: []string{"suggestions", "strengths", "weaknesses", "summary"},
}

const applicationEvaluationPromptBase = `You are an expert college admissions consultant. Review the following list of 
extracurricular activities to evaluate their depth, impact, and alignment with the student’s goals. Consider the 
following criteria to assess their quality and suggest improvements, ensuring strong activities are prioritized.
- Criteria:
  - Depth vs. Breadth:
    Does the student have a few activities with deep involvement (e.g., multiple years, significant roles) or many with 
    superficial involvement? Are there activities they have been involved in for multiple years (e.g., 2-4 years)?
  - Leadership and Impact:
    Has the student held leadership positions (e.g., captain, president, organizer)? What specific contributions 
    or achievements have they made (e.g., organizing events, mentoring others)? Are there awards, recognitions, or 
    outcomes that highlight their impact (e.g., team wins, community recognition)?
  - Relevance:
    Do the activities align with the student’s stated interests, goals, or intended major?
    Do the activities show a progression of involvement (e.g., starting as a member and becoming a leader)?
  - Order and Presentation:
    Are the most impressive or relevant activities listed first in the application?
    Are descriptions clear, concise, and impactful, highlighting specific achievements?

Extracurricular activities:
`
