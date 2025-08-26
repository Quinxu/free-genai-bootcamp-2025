package service

import (
	"database/sql"
	"fmt"

	"lang-portal/internal/models"
)

// StudyService handles study session related business logic
type StudyService struct {
	db *sql.DB
}

// NewStudyService creates a new StudyService
func NewStudyService(db *sql.DB) *StudyService {
	return &StudyService{db: db}
}

// StartStudySession starts a new study session
func (s *StudyService) StartStudySession(groupID, activityID int64) (*models.StudySession, error) {
	var session models.StudySession

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = tx.QueryRow(`
		INSERT INTO study_sessions (group_id, study_activity_id)
		VALUES (?, ?)
		RETURNING id, created_at
	`, groupID, activityID).Scan(&session.ID, &session.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to create study session: %w", err)
	}

	session.GroupID = groupID
	session.StudyActivityID = activityID

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &session, nil
}

// RecordWordReview records a word review in a study session
func (s *StudyService) RecordWordReview(sessionID, wordID int64, correct bool) (*models.WordReviewItem, error) {
	var review models.WordReviewItem

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = tx.QueryRow(`
		INSERT INTO word_review_items (word_id, study_session_id, correct)
		VALUES (?, ?, ?)
		RETURNING id, created_at
	`, wordID, sessionID, correct).Scan(&review.ID, &review.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to create word review: %w", err)
	}

	review.WordID = wordID
	review.StudySessionID = sessionID
	review.Correct = correct

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &review, nil
}

// GetStudyProgress returns study progress statistics
func (s *StudyService) GetStudyProgress() (*models.StudyProgress, error) {
	var progress models.StudyProgress

	err := s.db.QueryRow(`
		SELECT 
			COUNT(DISTINCT word_id) as total_words_studied,
			(SELECT COUNT(*) FROM words) as total_available_words
		FROM word_review_items
	`).Scan(&progress.TotalWordsStudied, &progress.TotalAvailableWords)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch study progress: %w", err)
	}

	return &progress, nil
}

// GetQuickStats returns quick study statistics
func (s *StudyService) GetQuickStats() (*models.QuickStats, error) {
	var stats models.QuickStats

	err := s.db.QueryRow(`
		WITH stats AS (
			SELECT 
				COALESCE(AVG(CASE WHEN correct THEN 1.0 ELSE 0.0 END) * 100, 0) as success_rate,
				COUNT(DISTINCT study_session_id) as total_sessions,
				COUNT(DISTINCT 
					CASE 
						WHEN created_at >= datetime('now', '-30 days')
						THEN study_session_id 
					END
				) as recent_sessions
			FROM word_review_items
		),
		active_groups AS (
			SELECT COUNT(DISTINCT group_id) as count
			FROM study_sessions
			WHERE created_at >= datetime('now', '-30 days')
		),
		streak AS (
			SELECT COUNT(DISTINCT date(created_at)) as days
			FROM study_sessions
			WHERE created_at >= (
				SELECT MAX(created_at)
				FROM (
					SELECT created_at,
						   julianday(created_at) - julianday(LAG(created_at) OVER (ORDER BY created_at)) as gap
					FROM study_sessions
					ORDER BY created_at DESC
				)
				WHERE gap > 1 OR gap IS NULL
			)
		)
		SELECT 
			stats.success_rate,
			stats.total_sessions,
			active_groups.count,
			streak.days
		FROM stats, active_groups, streak
	`).Scan(
		&stats.SuccessRate,
		&stats.TotalStudySessions,
		&stats.TotalActiveGroups,
		&stats.StudyStreakDays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch quick stats: %w", err)
	}

	return &stats, nil
}

// GetLastStudySession returns the most recent study session
func (s *StudyService) GetLastStudySession() (*models.StudySessionWithDetails, error) {
	var session models.StudySessionWithDetails

	err := s.db.QueryRow(`
		SELECT 
			ss.id,
			ss.group_id,
			ss.study_activity_id,
			ss.created_at,
			g.name as group_name,
			COUNT(wri.id) as reviewed_words
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.id = (
			SELECT id FROM study_sessions 
			ORDER BY created_at DESC 
			LIMIT 1
		)
		GROUP BY ss.id
	`).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.CreatedAt,
		&session.GroupName,
		&session.ReviewedWords,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch last study session: %w", err)
	}

	return &session, nil
}

// GetStudySessions returns paginated study sessions
func (s *StudyService) GetStudySessions(page, perPage int) (*models.PaginatedResponse[models.StudySessionSummary], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	rows, err := s.db.Query(`
        SELECT ss.id, ss.group_id, ss.study_activity_id, ss.created_at, g.name as group_name
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?
    `, perPage, (page-1)*perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch study sessions: %w", err)
	}
	defer rows.Close()

	var items []models.StudySessionSummary
	for rows.Next() {
		var item models.StudySessionSummary
		if err := rows.Scan(&item.ID, &item.GroupID, &item.StudyActivityID, &item.CreatedAt, &item.GroupName); err != nil {
			return nil, fmt.Errorf("failed to scan study session: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating study sessions: %w", err)
	}

	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count study sessions: %w", err)
	}

	return &models.PaginatedResponse[models.StudySessionSummary]{
		Items: items,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

// GetStudySession returns a single study session summary
func (s *StudyService) GetStudySession(id int64) (*models.StudySessionSummary, error) {
	var item models.StudySessionSummary
	err := s.db.QueryRow(`
        SELECT ss.id, ss.group_id, ss.study_activity_id, ss.created_at, g.name as group_name
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        WHERE ss.id = ?
    `, id).Scan(&item.ID, &item.GroupID, &item.StudyActivityID, &item.CreatedAt, &item.GroupName)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch study session: %w", err)
	}
	return &item, nil
}

// GetStudySessionWords returns words associated with a study session
func (s *StudyService) GetStudySessionWords(sessionID int64, page, perPage int) (*models.PaginatedResponse[models.Word], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	rows, err := s.db.Query(`
        SELECT w.id, w.chinese, w.english, w.parts, w.created_at
        FROM words w
        JOIN word_review_items wri ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
        GROUP BY w.id
        ORDER BY MAX(wri.created_at) DESC
        LIMIT ? OFFSET ?
    `, sessionID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch session words: %w", err)
	}
	defer rows.Close()

	var items []models.Word
	for rows.Next() {
		var w models.Word
		if err := rows.Scan(&w.ID, &w.Chinese, &w.English, &w.Parts, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		items = append(items, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session words: %w", err)
	}

	var total int
	if err := s.db.QueryRow(`
        SELECT COUNT(DISTINCT w.id)
        FROM words w
        JOIN word_review_items wri ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
    `, sessionID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count session words: %w", err)
	}

	return &models.PaginatedResponse[models.Word]{
		Items: items,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

// GetStudySessionsByActivity returns paginated sessions for a given activity
func (s *StudyService) GetStudySessionsByActivity(activityID int64, page, perPage int) (*models.PaginatedResponse[models.ActivitySessionListItem], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	rows, err := s.db.Query(`
        SELECT 
            ss.id,
            g.name as group_name,
            ss.created_at as start_time,
            MAX(wri.created_at) as end_time,
            COUNT(wri.id) as review_items_count
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
        WHERE ss.study_activity_id = ?
        GROUP BY ss.id
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?
    `, activityID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch activity sessions: %w", err)
	}
	defer rows.Close()

	var items []models.ActivitySessionListItem
	for rows.Next() {
		var item models.ActivitySessionListItem
		var startTime sql.NullString
		var endTime sql.NullString
		if err := rows.Scan(&item.ID, &item.GroupName, &startTime, &endTime, &item.ReviewItemsCount); err != nil {
			return nil, fmt.Errorf("failed to scan activity session: %w", err)
		}
		item.StartTime = startTime.String
		if endTime.Valid {
			item.EndTime = endTime.String
		}
		if details, err := s.GetStudyActivityDetails(activityID); err == nil {
			item.ActivityName = details.Name
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activity sessions: %w", err)
	}

	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM study_sessions WHERE study_activity_id = ?", activityID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count activity sessions: %w", err)
	}

	return &models.PaginatedResponse[models.ActivitySessionListItem]{
		Items: items,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

// GetStudyActivityDetails returns details for a study activity from a static catalog
func (s *StudyService) GetStudyActivityDetails(activityID int64) (*models.StudyActivityDetails, error) {
	catalog := map[int64]models.StudyActivityDetails{
		1: {
			Name:         "Vocabulary Quiz",
			ThumbnailURL: "https://example.com/thumbnail.jpg",
			Description:  "Practice your vocabulary with flashcards",
		},
	}
	if v, ok := catalog[activityID]; ok {
		v.ID = activityID
		return &v, nil
	}
	return nil, sql.ErrNoRows
}
