package models

// StudyActivity represents a specific study activity
type StudyActivity struct {
	Base
	StudySessionID int64 `json:"study_session_id" db:"study_session_id"`
	GroupID        int64 `json:"group_id" db:"group_id"`
}

// StudyActivityDetails includes additional presentation details
type StudyActivityDetails struct {
	StudyActivity
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}
