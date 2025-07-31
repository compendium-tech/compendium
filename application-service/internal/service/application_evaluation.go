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

type ApplicationEvaluationService interface {
	EvaluateCurrentApplication(ctx context.Context) (*domain.ApplicationEvaluationResponse, error)
}

type applicationEvaluationService struct {
	applicationRepository repository.ApplicationRepository
	llmService            service.LLMService
}

func NewApplicationAssessmentService(
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

	supplementalEssaysEvaluation, err := s.evaluateSupplementalEssays(ctx, supplementalEssays)
	if err != nil {
		return nil, err
	}

	return &domain.ApplicationEvaluationResponse{
		ActivitiesEvaluationResponse:         *activitiesEvaluation,
		HonorsEvaluationResponse:             *honorsEvaluation,
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
	}, &activitiesEvaluationSchema)
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
	}, &honorsEvaluationSchema)
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

func (s *applicationEvaluationService) evaluateSupplementalEssays(ctx context.Context, essays []model.SupplementalEssay) (*domain.SupplementalEssaysEvaluationResponse, error) {
	prompt := supplementalEssaysPromptBase

	for idx, essay := range essays {
		prompt += fmt.Sprintf("%d. %s - %s\n", idx+1, essay.Title)
		prompt += essay.Content + "\n\n\n"
	}

	structuredOutputSchema := generateSupplementalEssaysEvaluationSchema(len(essays))

	llmResponse, err := s.llmService.GenerateResponse(ctx, []llmdomain.Message{
		{
			Role: llmdomain.RoleSystem,
			Text: prompt,
		},
	}, &structuredOutputSchema)
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

const activitiesEvaluationPromptBase = `
You are an expert college admissions consultant. Review the following list of 
extracurricular activities to evaluate their depth, impact, and alignment with the student’s goals. Consider the 
following criteria to assess their quality and suggest improvements, ensuring strong activities are prioritized.

## Criteria
  - **Depth vs. Breadth**:
    Does the student have a few activities with deep involvement (e.g., multiple years, significant roles) or many with 
    superficial involvement? Are there activities they have been involved in for multiple years (e.g., 2-4 years)?
  - **Leadership and Impact**:
    Has the student held leadership positions (e.g., captain, president, organizer)? What specific contributions 
    or achievements have they made (e.g., organizing events, mentoring others)? Are there awards, recognitions, or 
    outcomes that highlight their impact (e.g., team wins, community recognition)?
  - **Relevance**:
    Do the activities align with the student’s stated interests, goals, or intended major?
    Do the activities show a progression of involvement (e.g., starting as a member and becoming a leader)?
  - **Order and Presentation**:
    Are the most impressive or relevant activities listed first in the application?
    Are descriptions clear, concise, and impactful, highlighting specific achievements?

## Extracurricular activities to evaluate
`

const honorsEvaluationPromptBase = `
You are an expert college admissions consultant. Review the following list of honors and awards to evaluate their prestige, relevance, and impact. Provide a detailed analysis based on the following criteria.

## Criteria
- What honors or awards has the student received, and at what level (school, regional, national, international)?
- Are these honors relevant to the student’s interests, goals, or intended major?
- Do they demonstrate exceptional achievement or recognition (e.g., scholarships, academic competitions)?
- Are there any gaps where honors might be expected but are missing (e.g., no academic honors despite strong grades)?
- Are the honors listed in order of prestige or impact (e.g., national awards before school awards)?

## Honors to evaluate
`

var supplementalEssaysPromptBase = `
You are an expert college admissions consultant. Review the following supplemental essays to evaluate their 
relevance, specificity, and quality. Ensure they demonstrate genuine interest in the college and avoid repetition 
with other sections. Provide a detailed analysis based on the following criteria.

## Criteria
- For each supplemental essay, what is the prompt, and how well does the student address it?
- Does the essay provide specific reasons why the student wants to attend the college (e.g., unique programs, faculty, 
campus culture, diversity)?
- Is there new information or a new perspective that isn’t covered in the main essays or other sections?
- Does the essay feel tailored to the college, or could it be submitted to any school?
- Is the writing clear, engaging, and free of errors?

## Essays to evaluate`

var activitiesEvaluationSchema = llmdomain.StructuredOutputSchema{
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

var honorsEvaluationSchema = llmdomain.StructuredOutputSchema{
	Type: llmdomain.TypeObject,
	Properties: map[string]llmdomain.StructuredOutputSchema{
		"suggestions": {
			Type:        llmdomain.TypeString,
			Description: `A list of actionable recommendations to improve the honors section, such as reordering honors by prestige, clarifying relevance to the student’s goals, or pursuing additional awards in their field of interest.`,
		},
		"strengths": {
			Type:        llmdomain.TypeString,
			Description: `A list of key strengths in the honors section, highlighting prestigious awards, relevance to the student’s interests or intended major, or exceptional achievements at various levels (school, regional, national, international).`,
		},
		"weaknesses": {
			Type:        llmdomain.TypeString,
			Description: `A list of weaknesses in the honors section, such as limited high-level awards, lack of relevance to the student’s goals, or missing honors in expected areas based on academic or extracurricular performance.`,
		},
		"summary": {
			Type:        llmdomain.TypeString,
			Description: `A concise summary of the overall quality of the honors section, addressing prestige, relevance, and impact, with an evaluation of how well it reflects the student’s achievements and alignment with their goals.`,
		},
	},
	Required: []string{"suggestions", "strengths", "weaknesses", "summary"},
}

var supplementalEssaysEvaluationSchema = llmdomain.StructuredOutputSchema{
	Type: llmdomain.TypeObject,
	Properties: map[string]llmdomain.StructuredOutputSchema{
		"suggestions": {
			Type:        llmdomain.TypeString,
			Description: `A list of actionable recommendations to improve the supplemental essays, such as adding specific details about the college, strengthening personal connections, avoiding overlap with other sections, or correcting writing errors.`,
		},
		"strengths": {
			Type:        llmdomain.TypeString,
			Description: `A list of key strengths in the supplemental essays, highlighting specific references to the college’s programs or culture, clear and engaging writing, or unique perspectives that demonstrate genuine interest.`,
		},
		"weaknesses": {
			Type:        llmdomain.TypeString,
			Description: `A list of weaknesses in the supplemental essays, such as generic content that could apply to any college, overlap with other application sections, lack of personal connection, or issues with writing clarity or errors.`,
		},
		"summary": {
			Type:        llmdomain.TypeString,
			Description: `A concise summary of the overall quality of the supplemental essays, addressing relevance, specificity, writing quality, and alignment with the college’s values, with an evaluation of how well they demonstrate the student’s fit and interest.`,
		},
	},
	Required: []string{"suggestions", "strengths", "weaknesses", "summary"},
}

func generateSupplementalEssaysEvaluationSchema(essaysCount int) llmdomain.StructuredOutputSchema {
	essaysCountInt64 := int64(essaysCount)

	return llmdomain.StructuredOutputSchema{
		Type: llmdomain.TypeObject,
		Properties: map[string]llmdomain.StructuredOutputSchema{
			"individualEvaluations": {
				Type:     llmdomain.TypeArray,
				MinItems: &essaysCountInt64,
				MaxItems: &essaysCountInt64,
				Items: &llmdomain.StructuredOutputSchema{
					Type: llmdomain.TypeObject,
					Properties: map[string]llmdomain.StructuredOutputSchema{
						"suggestions": {
							Type: llmdomain.TypeArray,
							Items: &llmdomain.StructuredOutputSchema{
								Type:        llmdomain.TypeString,
								Description: `An actionable recommendation to improve the specific supplemental essay, such as adding specific details about the college, strengthening personal connections, avoiding overlap with other sections, or correcting writing errors.`,
							},
							Description: `A list of actionable recommendations to improve the specific supplemental essay.`,
						},
						"strengths": {
							Type: llmdomain.TypeArray,
							Items: &llmdomain.StructuredOutputSchema{
								Type:        llmdomain.TypeString,
								Description: `A key strength of the specific supplemental essay, highlighting specific references to the college’s programs or culture, clear and engaging writing, or unique perspectives that demonstrate genuine interest.`,
							},
							Description: `A list of key strengths of the specific supplemental essay.`,
						},
						"weaknesses": {
							Type: llmdomain.TypeArray,
							Items: &llmdomain.StructuredOutputSchema{
								Type:        llmdomain.TypeString,
								Description: `A weakness of the specific supplemental essay, such as generic content that could apply to any college, overlap with other application sections, lack of personal connection, or issues with writing clarity or errors.`,
							},
							Description: `A list of weaknesses of the specific supplemental essay.`,
						},
					},
					Required: []string{"suggestions", "strengths", "weaknesses"},
				},
				Description: `A list of evaluations, one for each supplemental essay, where each evaluation includes suggestions, strengths, and weaknesses for the specific essay. The order of evaluations must match the order of the essays (e.g., the first evaluation corresponds to the first essay, the second to the second essay, etc.).`,
			},
			"assessment": {
				Type:        llmdomain.TypeString,
				Description: `A concise summary of the overall quality of all supplemental essays, addressing relevance, specificity, writing quality, and alignment with the college’s values, with an evaluation of how well they collectively demonstrate the student’s fit and interest.`,
			},
		},
		Required: []string{"individualEvaluations", "assessment"},
	}
}
