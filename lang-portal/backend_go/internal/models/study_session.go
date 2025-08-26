package models

// StudySession represents a study session
type StudySession struct {
	Base
	GroupID         int64 `json:"group_id" db:"group_id"`
	StudyActivityID int64 `json:"study_activity_id" db:"study_activity_id"`
}

// StudySessionWithDetails includes additional details about the study session
type StudySessionWithDetails struct {
	StudySession
	GroupName     string `json:"group_name"`
	ActivityName  string `json:"activity_name"`
	ReviewedWords int    `json:"reviewed_words"`
}

// StudySessionSummary represents the list/detail shape required by the API spec
type StudySessionSummary struct {
	Base
	GroupID         int64  `json:"group_id"`
	StudyActivityID int64  `json:"study_activity_id"`
	GroupName       string `json:"group_name"`
}

// ActivitySessionListItem matches spec for study activity sessions list
type ActivitySessionListItem struct {
	ID               int64  `json:"id"`
	ActivityName     string `json:"activity_name"`
	GroupName        string `json:"group_name"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	ReviewItemsCount int    `json:"review_items_count"`
}
