package domain

type ApplicationEvaluationResponse struct {
	ActivitiesEvaluationResponse `json:"activities_evaluation"`
}

type ActivitiesEvaluationResponse struct {
	Suggestions []string           `json:"suggestions"`
	Strengths   []string           `json:"strengths"`
	Weaknesses  []string           `json:"weaknesses"`
	Summary     string             `json:"assessment"`
	Activities  []ActivityResponse `json:"activities"`
}
