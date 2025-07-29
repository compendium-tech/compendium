package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
}

type Activity struct {
	ApplicationID uuid.UUID
	Name          string
	Role          string
	Description   *string
	HoursPerWeek  int
	WeeksPerYear  int
	Category      ActivityCategory
	Grades        []Grade
}

type Honor struct {
	ApplicationID uuid.UUID
	Title         string
	Description   *string
	Level         HonorLevel
	Grade         Grade
}

type ActivityCategory string

const (
	ActivityCategoryAcademic               ActivityCategory = "academic"
	ActivityCategoryArt                    ActivityCategory = "art"
	ActivityCategoryAthletics              ActivityCategory = "athletics"
	ActivityCategoryCareerOriented         ActivityCategory = "career_oriented"
	ActivityCategoryCommunityService       ActivityCategory = "community_service"
	ActivityCategoryCultural               ActivityCategory = "cultural"
	ActivityCategoryDebateSpeech           ActivityCategory = "debate_speech"
	ActivityCategoryEnvironmental          ActivityCategory = "environmental"
	ActivityCategoryFamilyResponsibilities ActivityCategory = "family_responsibilities"
	ActivityCategoryJournalismPublication  ActivityCategory = "journalism_publication"
	ActivityCategoryMusic                  ActivityCategory = "music"
	ActivityCategoryReligious              ActivityCategory = "religious"
	ActivityCategoryResearch               ActivityCategory = "research"
	ActivityCategoryRobotics               ActivityCategory = "robotics"
	ActivityCategorySchoolSpirit           ActivityCategory = "school_spirit"
	ActivityCategoryStudentGovernment      ActivityCategory = "student_government"
	ActivityCategoryTheatreDrama           ActivityCategory = "theatre_drama"
	ActivityCategoryWork                   ActivityCategory = "work"
	ActivityCategoryOther                  ActivityCategory = "other"
)

func (a *ActivityCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(ActivityCategoryAcademic),
		string(ActivityCategoryArt),
		string(ActivityCategoryAthletics),
		string(ActivityCategoryCareerOriented),
		string(ActivityCategoryCommunityService),
		string(ActivityCategoryCultural),
		string(ActivityCategoryDebateSpeech),
		string(ActivityCategoryEnvironmental),
		string(ActivityCategoryFamilyResponsibilities),
		string(ActivityCategoryJournalismPublication),
		string(ActivityCategoryMusic),
		string(ActivityCategoryReligious),
		string(ActivityCategoryResearch),
		string(ActivityCategoryRobotics),
		string(ActivityCategorySchoolSpirit),
		string(ActivityCategoryStudentGovernment),
		string(ActivityCategoryTheatreDrama),
		string(ActivityCategoryWork),
		string(ActivityCategoryOther):
		*a = ActivityCategory(s)
		return nil
	}
	return fmt.Errorf("invalid activity category: %s", s)
}

type Grade string

const (
	Grade9            Grade = "9"
	Grade10           Grade = "10"
	Grade11           Grade = "11"
	Grade12           Grade = "12"
	GradePostGraduate Grade = "post_graduate"
)

func (g *Grade) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(Grade9),
		string(Grade10),
		string(Grade11),
		string(Grade12),
		string(GradePostGraduate):
		*g = Grade(s)
		return nil
	}
	return fmt.Errorf("invalid grade: %s", s)
}

type HonorLevel string

const (
	HonorLevelSchool        HonorLevel = "school"
	HonorLevelRegional      HonorLevel = "regional"
	HonorLevelNational      HonorLevel = "national"
	HonorLevelInternational HonorLevel = "international"
)

func (h *HonorLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(HonorLevelSchool),
		string(HonorLevelRegional),
		string(HonorLevelNational),
		string(HonorLevelInternational):
		*h = HonorLevel(s)
		return nil
	}
	return fmt.Errorf("invalid honor level: %s", s)
}
