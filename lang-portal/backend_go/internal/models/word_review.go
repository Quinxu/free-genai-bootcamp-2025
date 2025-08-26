package models

// WordReviewItem represents a word review record
type WordReviewItem struct {
	Base
	WordID          int64 `json:"word_id" db:"word_id"`
	StudySessionID  int64 `json:"study_session_id" db:"study_session_id"`
	Correct         bool  `json:"correct" db:"correct"`
}

// WordReviewStats represents statistics for word reviews
type WordReviewStats struct {
	TotalReviews  int     `json:"total_reviews"`
	CorrectCount  int     `json:"correct_count"`
	WrongCount    int     `json:"wrong_count"`
	SuccessRate   float64 `json:"success_rate"`
}
