package domain

type ApplicationEvaluationResponse struct {
	ActivitiesEvaluationResponse         `json:"activitiesEvaluation"`
	HonorsEvaluationResponse             `json:"honorsEvaluation"`
	SupplementalEssaysEvaluationResponse `json:"supplementalEssaysEvaluation"`
}

type ActivitiesEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Summary     string   `json:"assessment"`
}

type HonorsEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Summary     string   `json:"assessment"`
}

type SupplementalEssaysEvaluationResponse struct {
	IndividualEvaluations []SupplementalEssayEvaluationResponse `json:"individualEvaluations"`
	Summary               string                                `json:"assessment"`
}

type SupplementalEssayEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
}
