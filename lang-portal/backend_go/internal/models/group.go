package models

// Group represents a thematic group of words
type Group struct {
	Base
	Name string `json:"name" db:"name"`
}

// GroupWithWords represents a group with its associated words
type GroupWithWords struct {
	Group
	Words []Word `json:"words"`
}

// GroupStats represents statistics for a group
type GroupStats struct {
	TotalWords     int     `json:"total_words"`
	StudiedWords   int     `json:"studied_words"`
	SuccessRate    float64 `json:"success_rate"`
	LastStudiedAt  string  `json:"last_studied_at,omitempty"`
}
