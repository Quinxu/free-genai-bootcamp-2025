package service

import (
	"database/sql"
	"fmt"

	"lang-portal/internal/database/query"
	"lang-portal/internal/models"
)

// GroupService handles group-related business logic
type GroupService struct {
	db *sql.DB
}

// NewGroupService creates a new GroupService
func NewGroupService(db *sql.DB) *GroupService {
	return &GroupService{db: db}
}

// GetGroups returns a paginated list of groups
func (s *GroupService) GetGroups(page, perPage int) (*models.PaginatedResponse[models.Group], error) {
	q := query.New("SELECT * FROM groups")
	q.OrderBy("name ASC").Paginate(page, perPage)

	rows, err := q.Execute(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating groups: %w", err)
	}

	total, err := q.ExecuteCount(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count groups: %w", err)
	}

	return &models.PaginatedResponse[models.Group]{
		Items: groups,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

// GetGroupByID returns a single group with its words
func (s *GroupService) GetGroupByID(id int64) (*models.GroupWithWords, error) {
	// First get the group
	var group models.Group
	err := s.db.QueryRow("SELECT * FROM groups WHERE id = ?", id).
		Scan(&group.ID, &group.Name, &group.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group: %w", err)
	}

	// Then get its words
	rows, err := s.db.Query(`
		SELECT w.* FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.created_at DESC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group words: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var w models.Word
		var parts []byte
		if err := rows.Scan(&w.ID, &w.Chinese, &w.English, &parts, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating words: %w", err)
	}

	return &models.GroupWithWords{
		Group: group,
		Words: words,
	}, nil
}

// GetGroupWordsPaginated returns paginated words for a group
func (s *GroupService) GetGroupWordsPaginated(groupID int64, page, perPage int) (*models.PaginatedResponse[models.Word], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	rows, err := s.db.Query(`
        SELECT w.id, w.chinese, w.english, w.parts, w.created_at
        FROM words w
        JOIN words_groups wg ON w.id = wg.word_id
        WHERE wg.group_id = ?
        ORDER BY w.created_at DESC
        LIMIT ? OFFSET ?
    `, groupID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group words: %w", err)
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
		return nil, fmt.Errorf("error iterating group words: %w", err)
	}

	var total int
	if err := s.db.QueryRow(`
        SELECT COUNT(*)
        FROM words_groups
        WHERE group_id = ?
    `, groupID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count group words: %w", err)
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

// GetGroupStudySessions returns paginated study sessions for a group
func (s *GroupService) GetGroupStudySessions(groupID int64, page, perPage int) (*models.PaginatedResponse[models.StudySessionSummary], error) {
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
        WHERE ss.group_id = ?
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?
    `, groupID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group sessions: %w", err)
	}
	defer rows.Close()

	var items []models.StudySessionSummary
	for rows.Next() {
		var item models.StudySessionSummary
		if err := rows.Scan(&item.ID, &item.GroupID, &item.StudyActivityID, &item.CreatedAt, &item.GroupName); err != nil {
			return nil, fmt.Errorf("failed to scan group session: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group sessions: %w", err)
	}

	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM study_sessions WHERE group_id = ?", groupID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count group sessions: %w", err)
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

// GetGroupStats returns statistics for a group
func (s *GroupService) GetGroupStats(id int64) (*models.GroupStats, error) {
	var stats models.GroupStats

	err := s.db.QueryRow(`
		WITH group_stats AS (
			SELECT 
				COUNT(DISTINCT w.id) as total_words,
				COUNT(DISTINCT CASE WHEN wri.id IS NOT NULL THEN w.id END) as studied_words,
				COALESCE(AVG(CASE WHEN wri.correct THEN 1.0 ELSE 0.0 END), 0) as success_rate,
				MAX(wri.created_at) as last_studied_at
			FROM words w
			JOIN words_groups wg ON w.id = wg.word_id
			LEFT JOIN word_review_items wri ON w.id = wri.word_id
			WHERE wg.group_id = ?
		)
		SELECT 
			total_words,
			studied_words,
			success_rate * 100,
			last_studied_at
		FROM group_stats
	`, id).Scan(&stats.TotalWords, &stats.StudiedWords, &stats.SuccessRate, &stats.LastStudiedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch group stats: %w", err)
	}

	return &stats, nil
}
