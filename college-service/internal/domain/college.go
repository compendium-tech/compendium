package domain

type SearchCollegesRequest struct {
	PageIndex          *int    `json:"pageIndex"`
	StateOrCountry     *string `json:"stateOrCountry"`
	SemanticSearchText *string `json:"semanticSearchText" validate:"max=1000"`
}

type CollegeResponse struct {
	Name           string `json:"name"`
	City           string `json:"city"`
	StateOrCountry string `json:"stateOrCountry"`
	Description    string `json:"description"`
}
