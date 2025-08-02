package service

import "github.com/compendium-tech/compendium/application-service/internal/domain"

const applicationEvaluationPromptBase = `
You are an expert college admissions consultant. Evaluate the entire college application, including academics,
character, extracurricular activities, essays (personal statement, teacher recommendations, counselor recommendation),
honors, supplemental essays, and authenticity/fit with the target college. Provide a detailed analysis for each
section based on the specified criteria, ensuring each section is evaluated distinctly. Identify any overlaps between
sections (e.g., essays repeating activities or honors) to avoid redundancy. Synthesize the assessments into a
cohesive overall picture of the student, highlighting their strengths, weaknesses, and alignment with the
college’s expectations. The number of evaluations for each section must match the number of items provided
(e.g., one evaluation per essay type, one for academics, etc.), and the order of evaluations must correspond
to the order of the input items.

# Criteria

## Character

- Does the application present a clear, consistent picture of the student’s character (e.g., resilience, empathy, leadership)?
- Are there specific examples of positive traits (e.g., integrity, perseverance) across essays, activities, or recommendations?
- How does character come through in the personal statement, activities, and recommendations?
- Are there inconsistencies or gaps raising questions (e.g., unexplained activity gaps, conflicting narratives)?
- Does the application reflect authenticity and self-awareness?

## Extracurricular Activities

### Depth vs. Breadth:
- Does the student have deep involvement in a few activities (e.g., multiple years, significant roles) or superficial involvement in many?
- Are there activities with multi-year commitment (2-4 years)?

### Leadership and Impact:
- Has the student held leadership roles (e.g., captain, president)?
- What specific contributions or achievements are highlighted (e.g., organizing events)?
- Are there awards or outcomes showing impact (e.g., team wins, community recognition)?

### Relevance:
- Do activities align with the student’s interests, goals, or intended major?
- Do they show progression (e.g., member to leader)?

### Order and Presentation:
- Are the most impressive activities listed first?
- Are descriptions clear, concise, and impactful?

## Essays (Personal Statement, Teacher Recommendations, Counselor Recommendation)

### Personal Statement:
- What is the main theme or story?
- Does it provide new insights into the student’s background or aspirations?
- Are there overused or clichéd topics?
- Does it avoid restating activities or honors?
- Is the writing clear, engaging, and free of errors?
- Does it feel genuine, with a unique voice?

### Teacher Recommendations:
- Do they provide specific examples of strengths (e.g., academic curiosity)?
- Are there details beyond general praise?
- Do they convey enthusiasm and knowledge of the student?
- Do they align with and complement the application without repetition?

### Counselor Recommendation:
- Do they provide specific examples of character or contributions?
- Are there details beyond general praise?
- Do they convey enthusiasm and knowledge of the student in a school context?
- Do they align with and complement the application?

## Honors

- What honors are received, and at what level (school, regional, national, international)?
- Are they relevant to the student’s interests or major?
- Do they demonstrate exceptional achievement (e.g., scholarships, competitions)?
- Are there gaps where honors are expected but missing?
- Are honors listed in order of prestige?

## Supplemental Essays

- What is the prompt, and how well is it addressed?
- Do they provide specific reasons for wanting to attend the college (e.g., programs, faculty)?
- Do they offer new information not covered elsewhere?
- Are they tailored to the college, or generic?
- Is the writing clear, engaging, and error-free?

## Authenticity and Fit

- Does the application show genuine interest in the college (e.g., specific programs, values)?
- Are there inconsistencies raising questions (e.g., essays mentioning passions not in activities)?
- Does the student demonstrate clear goals and how the college supports them?
- Is there evidence of demonstrated interest (e.g., campus visits) if required?
- Does the application reflect an authentic voice, or is it overly polished?
`

func generateApplicationEvaluationSchema(essaysCount int, supplementalEssaysCount int) domain.LLMSchema {
	return domain.LLMSchema{
		Type: domain.TypeObject,
		Properties: map[string]domain.LLMSchema{
			"suggestions": {
				Type:        domain.TypeString,
				Description: `A list of actionable recommendations to improve the overall application, synthesizing suggestions across all sections (academics, character, activities, essays, honors, supplemental essays, interview, and authenticity/fit) to enhance cohesiveness, alignment with the college’s expectations, or address gaps and weaknesses.`,
			},
			"strengths": {
				Type:        domain.TypeString,
				Description: `A list of key strengths across the entire application, highlighting standout qualities such as exceptional academic performance, leadership, compelling narratives, or strong alignment with the college’s values and programs.`,
			},
			"weaknesses": {
				Type:        domain.TypeString,
				Description: `A list of weaknesses across the entire application, identifying areas such as inconsistencies, lack of depth, overlap between sections, or misalignment with the student’s goals or the college’s expectations.`,
			},
			"summary": {
				Type:        domain.TypeString,
				Description: `A concise summary of the overall quality of the application, synthesizing the cohesiveness, strengths, weaknesses, and alignment with the college’s culture and expectations, presenting a holistic picture of the student’s character, achievements, and fit.`,
			},
			"activitiesEvaluation":         activitiesEvaluationSchema,
			"honorsEvaluation":             honorsEvaluationSchema,
			"essaysEvaluation":             generateEssaysEvaluationSchema(essaysCount),
			"supplementalEssaysEvaluation": generateSupplementalEssaysEvaluationSchema(supplementalEssaysCount),
		},
	}
}

var activitiesEvaluationSchema = domain.LLMSchema{
	Type: domain.TypeObject,
	Properties: map[string]domain.LLMSchema{
		"suggestions": {
			Type: domain.TypeString,
			Description: `A list of actionable recommendations to improve the activities section, such
as reordering activities, adding specific achievements, or increasing involvement in relevant activities.`,
		},
		"strengths": {
			Type: domain.TypeString,
			Description: `A list of key strengths in the activities section, highlighting deep involvement,
leadership roles, impactful contributions, or alignment with the student’s goals.`,
		},
		"weaknesses": {
			Type: domain.TypeString,
			Description: `A list of weaknesses in the activities section, such as superficial involvement, lack
of leadership, unclear descriptions, or misalignment with the student’s goals.`,
		},
		"summary": {
			Type: domain.TypeString,
			Description: `A concise summary of the overall quality of the activities section, addressing depth,
impact, relevance, and presentation, with an evaluation of how well it reflects the student’s strengths and goals.`,
		},
	},
	Required: []string{"suggestions", "strengths", "weaknesses", "summary"},
}

var honorsEvaluationSchema = domain.LLMSchema{
	Type: domain.TypeObject,
	Properties: map[string]domain.LLMSchema{
		"suggestions": {
			Type:        domain.TypeString,
			Description: `A list of actionable recommendations to improve the honors section, such as reordering honors by prestige, clarifying relevance to the student’s goals, or pursuing additional awards in their field of interest.`,
		},
		"strengths": {
			Type:        domain.TypeString,
			Description: `A list of key strengths in the honors section, highlighting prestigious awards, relevance to the student’s interests or intended major, or exceptional achievements at various levels (school, regional, national, international).`,
		},
		"weaknesses": {
			Type:        domain.TypeString,
			Description: `A list of weaknesses in the honors section, such as limited high-level awards, lack of relevance to the student’s goals, or missing honors in expected areas based on academic or extracurricular performance.`,
		},
		"summary": {
			Type:        domain.TypeString,
			Description: `A concise summary of the overall quality of the honors section, addressing prestige, relevance, and impact, with an evaluation of how well it reflects the student’s achievements and alignment with their goals.`,
		},
	},
	Required: []string{"suggestions", "strengths", "weaknesses", "summary"},
}

func generateEssaysEvaluationSchema(essaysCount int) domain.LLMSchema {
	essaysCountInt64 := int64(essaysCount)

	return domain.LLMSchema{
		Type: domain.TypeObject,
		Properties: map[string]domain.LLMSchema{
			"individualEvaluations": {
				Type:     domain.TypeArray,
				MinItems: &essaysCountInt64,
				MaxItems: &essaysCountInt64,
				Items: &domain.LLMSchema{
					Type: domain.TypeObject,
					Properties: map[string]domain.LLMSchema{
						"suggestions": {
							Type: domain.TypeArray,
							Items: &domain.LLMSchema{
								Type:        domain.TypeString,
								Description: `An actionable recommendation to improve the specific essay. For personal statements, this may include focusing on personal growth, avoiding clichés, or correcting writing errors. For recommendations, this may include adding specific anecdotes, highlighting diverse qualities, or aligning with the application.`,
							},
							Description: `A list of actionable recommendations to improve the specific essay (personal statement, teacher recommendation, or counselor recommendation).`,
						},
						"strengths": {
							Type: domain.TypeArray,
							Items: &domain.LLMSchema{
								Type:        domain.TypeString,
								Description: `A key strength of the specific essay. For personal statements, this may highlight compelling narratives, unique voice, or alignment with aspirations. For recommendations, this may highlight specific examples, enthusiastic tone, or alignment with the student’s application.`,
							},
							Description: `A list of key strengths of the specific essay (personal statement, teacher recommendation, or counselor recommendation).`,
						},
						"weaknesses": {
							Type: domain.TypeArray,
							Items: &domain.LLMSchema{
								Type:        domain.TypeString,
								Description: `A weakness of the specific essay. For personal statements, this may include overlap with other sections, lack of authenticity, or writing errors. For recommendations, this may include lack of specificity, generic praise, or misalignment with the application.`,
							},
							Description: `A list of weaknesses of the specific essay (personal statement, teacher recommendation, or counselor recommendation).`,
						},
					},
					Required: []string{"suggestions", "strengths", "weaknesses"},
				},
				Description: `A list of evaluations, one for each essay (personal statement, teacher recommendation, or counselor recommendation), where each evaluation includes suggestions, strengths, and weaknesses. The number of evaluations must equal the number of essays provided, and the order of evaluations must match the order of the essays (e.g., the first evaluation corresponds to the first essay, the second to the second essay, etc.).`,
			},
			"assessment": {
				Type:        domain.TypeString,
				Description: `A concise summary of the overall quality of all essays (personal statement, teacher recommendations, and counselor recommendation), addressing content, writing quality, authenticity, alignment with the student’s application, and how well they collectively convey the student’s character and fit.`,
			},
			"overlap": {
				Type: domain.TypeArray,
				Items: &domain.LLMSchema{
					Type:        domain.TypeString,
					Description: `A description of any overlap between an essay and other application sections (e.g., activities or honors), such as repeating specific activities or achievements that dilute the application’s impact.`,
				},
				Description: `A list of descriptions identifying any overlap between the essays and other application sections (e.g., activities or honors).`,
			},
		},
		Required: []string{"individualEvaluations", "assessment", "overlap"},
	}
}

func generateSupplementalEssaysEvaluationSchema(essaysCount int) domain.LLMSchema {
	essaysCountInt64 := int64(essaysCount)

	return domain.LLMSchema{
		Type: domain.TypeObject,
		Properties: map[string]domain.LLMSchema{
			"individualEvaluations": {
				Type:     domain.TypeArray,
				MinItems: &essaysCountInt64,
				MaxItems: &essaysCountInt64,
				Items: &domain.LLMSchema{
					Type: domain.TypeObject,
					Properties: map[string]domain.LLMSchema{
						"suggestions": {
							Type: domain.TypeArray,
							Items: &domain.LLMSchema{
								Type:        domain.TypeString,
								Description: `An actionable recommendation to improve the specific supplemental essay, such as adding specific details about the college, strengthening personal connections, avoiding overlap with other sections, or correcting writing errors.`,
							},
							Description: `A list of actionable recommendations to improve the specific supplemental essay.`,
						},
						"strengths": {
							Type: domain.TypeArray,
							Items: &domain.LLMSchema{
								Type:        domain.TypeString,
								Description: `A key strength of the specific supplemental essay, highlighting specific references to the college’s programs or culture, clear and engaging writing, or unique perspectives that demonstrate genuine interest.`,
							},
							Description: `A list of key strengths of the specific supplemental essay.`,
						},
						"weaknesses": {
							Type: domain.TypeArray,
							Items: &domain.LLMSchema{
								Type:        domain.TypeString,
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
				Type:        domain.TypeString,
				Description: `A concise summary of the overall quality of all supplemental essays, addressing relevance, specificity, writing quality, and alignment with the college’s values, with an evaluation of how well they collectively demonstrate the student’s fit and interest.`,
			},
		},
		Required: []string{"individualEvaluations", "assessment"},
	}
}
