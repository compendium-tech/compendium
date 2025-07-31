package domain

type ApplicationEvaluationResponse struct {
	ActivitiesEvaluationResponse         `json:"activitiesEvaluation"`
	HonorsEvaluationResponse             `json:"honorsEvaluation"`
	EssaysEvaluationResponse             `json:"essaysEvaluation"`
	SupplementalEssaysEvaluationResponse `json:"supplementalEssaysEvaluation"`
}

type ActivitiesEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Summary     string   `json:"summary"`
}

type HonorsEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Summary     string   `json:"summary"`
}

type EssaysEvaluationResponse struct {
	IndividualEvaluations []SupplementalEssayEvaluationResponse `json:"individualEvaluations"`
	Summary               string                                `json:"summary"`
}

type EssayEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
}

type SupplementalEssaysEvaluationResponse struct {
	IndividualEvaluations []SupplementalEssayEvaluationResponse `json:"individualEvaluations"`
	Summary               string                                `json:"summary"`
}

type SupplementalEssayEvaluationResponse struct {
	Suggestions []string `json:"suggestions"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
}
