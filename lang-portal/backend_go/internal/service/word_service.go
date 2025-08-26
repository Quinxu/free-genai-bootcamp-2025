package service

import (
	"database/sql"
	"fmt"

	"lang-portal/internal/database/query"
	"lang-portal/internal/models"
)

// WordService handles word-related business logic
type WordService struct {
	db *sql.DB
}

// NewWordService creates a new WordService
func NewWordService(db *sql.DB) *WordService {
	return &WordService{db: db}
}

// GetWords returns a paginated list of words with their stats
func (s *WordService) GetWords(page, perPage int) (*models.PaginatedResponse[models.WordWithStats], error) {
	q := query.New("SELECT w.*, " +
		"(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 1) as correct_count, " +
		"(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 0) as wrong_count " +
		"FROM words w")

	q.OrderBy("w.created_at DESC").Paginate(page, perPage)

	rows, err := q.Execute(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch words: %w", err)
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var w models.WordWithStats
		var parts []byte
		var correctCount, wrongCount int

		err := rows.Scan(
			&w.ID,
			&w.Chinese,
			&w.English,
			&parts,
			&w.CreatedAt,
			&correctCount,
			&wrongCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}

		// Preserve raw JSON for parts without unmarshalling
		w.Parts = parts

		w.Stats = models.WordStats{
			CorrectCount: correctCount,
			WrongCount:   wrongCount,
		}

		words = append(words, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating words: %w", err)
	}

	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count words: %w", err)
	}

	return &models.PaginatedResponse[models.WordWithStats]{
		Items: words,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

// GetWordByID returns a single word with its stats
func (s *WordService) GetWordByID(id int64) (*models.WordWithStats, error) {
	q := query.New("SELECT w.*, " +
		"(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 1) as correct_count, " +
		"(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 0) as wrong_count " +
		"FROM words w")
	q.Where("w.id = ?", id)

	rows, err := q.Execute(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch word: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var w models.WordWithStats
	var parts []byte
	var correctCount, wrongCount int

	err = rows.Scan(
		&w.ID,
		&w.Chinese,
		&w.English,
		&parts,
		&w.CreatedAt,
		&correctCount,
		&wrongCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan word: %w", err)
	}

	// Preserve raw JSON for parts
	w.Parts = parts

	w.Stats = models.WordStats{
		CorrectCount: correctCount,
		WrongCount:   wrongCount,
	}

	return &w, nil
}

// GetGroupsForWord returns the groups that contain a given word
func (s *WordService) GetGroupsForWord(wordID int64) ([]models.Group, error) {
	rows, err := s.db.Query(`
        SELECT g.id, g.name, g.created_at
        FROM groups g
        JOIN words_groups wg ON wg.group_id = g.id
        WHERE wg.word_id = ?
        ORDER BY g.name ASC
    `, wordID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch word groups: %w", err)
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
	return groups, nil
}
