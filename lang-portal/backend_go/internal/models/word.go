package models

import (
	"encoding/json"
)

// Word represents a vocabulary word in the system
type Word struct {
	Base
	Chinese string          `json:"chinese" db:"chinese"`
	English string          `json:"english" db:"english"`
	Parts   json.RawMessage `json:"parts" db:"parts"`
}

// WordStats represents statistics for a word
type WordStats struct {
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

// WordWithStats combines Word with its statistics
type WordWithStats struct {
	Word
	Stats WordStats `json:"stats"`
}
