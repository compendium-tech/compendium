package service

import llmdomain "github.com/compendium-tech/compendium/llm-common/pkg/domain"

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

func generateEssaysEvaluationSchema(essaysCount int) llmdomain.StructuredOutputSchema {
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
								Description: `An actionable recommendation to improve the specific essay. For personal statements, this may include focusing on personal growth, avoiding clichés, or correcting writing errors. For recommendations, this may include adding specific anecdotes, highlighting diverse qualities, or aligning with the application.`,
							},
							Description: `A list of actionable recommendations to improve the specific essay (personal statement, teacher recommendation, or counselor recommendation).`,
						},
						"strengths": {
							Type: llmdomain.TypeArray,
							Items: &llmdomain.StructuredOutputSchema{
								Type:        llmdomain.TypeString,
								Description: `A key strength of the specific essay. For personal statements, this may highlight compelling narratives, unique voice, or alignment with aspirations. For recommendations, this may highlight specific examples, enthusiastic tone, or alignment with the student’s application.`,
							},
							Description: `A list of key strengths of the specific essay (personal statement, teacher recommendation, or counselor recommendation).`,
						},
						"weaknesses": {
							Type: llmdomain.TypeArray,
							Items: &llmdomain.StructuredOutputSchema{
								Type:        llmdomain.TypeString,
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
				Type:        llmdomain.TypeString,
				Description: `A concise summary of the overall quality of all essays (personal statement, teacher recommendations, and counselor recommendation), addressing content, writing quality, authenticity, alignment with the student’s application, and how well they collectively convey the student’s character and fit.`,
			},
			"overlap": {
				Type: llmdomain.TypeArray,
				Items: &llmdomain.StructuredOutputSchema{
					Type:        llmdomain.TypeString,
					Description: `A description of any overlap between an essay and other application sections (e.g., activities or honors), such as repeating specific activities or achievements that dilute the application’s impact.`,
				},
				Description: `A list of descriptions identifying any overlap between the essays and other application sections (e.g., activities or honors).`,
			},
		},
		Required: []string{"individualEvaluations", "assessment", "overlap"},
	}
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
