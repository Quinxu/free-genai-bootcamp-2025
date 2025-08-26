package seeds

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Seeder handles database seeding
type Seeder struct {
	db   *sql.DB
	path string
}

// NewSeeder creates a new seeder
func NewSeeder(db *sql.DB, seedsPath string) *Seeder {
	return &Seeder{
		db:   db,
		path: seedsPath,
	}
}

// SeedData represents seed data configuration
type SeedData struct {
	Groups []struct {
		Name  string `json:"name"`
		Words []struct {
			Chinese string                 `json:"chinese"`
			English string                 `json:"english"`
			Parts   map[string]interface{} `json:"parts"`
		} `json:"words"`
	} `json:"groups"`
}

// Seed runs the database seeding process
func (s *Seeder) Seed() error {
	files, err := ioutil.ReadDir(s.path)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(s.path, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %w", file.Name(), err)
		}

		var data SeedData
		if err := json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to parse seed file %s: %w", file.Name(), err)
		}

		if err := s.seedGroups(data.Groups); err != nil {
			return fmt.Errorf("failed to seed groups from %s: %w", file.Name(), err)
		}
	}

	return nil
}

func (s *Seeder) seedGroups(groups []struct {
	Name  string `json:"name"`
	Words []struct {
		Chinese string                 `json:"chinese"`
		English string                 `json:"english"`
		Parts   map[string]interface{} `json:"parts"`
	} `json:"words"`
}) error {
	for _, group := range groups {
		// Insert group
		var groupID int64
		err := s.db.QueryRow(
			"INSERT INTO groups (name) VALUES (?) RETURNING id",
			group.Name,
		).Scan(&groupID)
		if err != nil {
			return fmt.Errorf("failed to insert group %s: %w", group.Name, err)
		}

		// Insert words and create word-group relationships
		for _, word := range group.Words {
			parts, err := json.Marshal(word.Parts)
			if err != nil {
				return fmt.Errorf("failed to marshal parts for word %s: %w", word.Chinese, err)
			}

			var wordID int64
			err = s.db.QueryRow(
				"INSERT INTO words (chinese, english, parts) VALUES (?, ?, ?) RETURNING id",
				word.Chinese, word.English, parts,
			).Scan(&wordID)
			if err != nil {
				return fmt.Errorf("failed to insert word %s: %w", word.Chinese, err)
			}

			_, err = s.db.Exec(
				"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
				wordID, groupID,
			)
			if err != nil {
				return fmt.Errorf("failed to create word-group relationship for word %s: %w", word.Chinese, err)
			}
		}
	}

	return nil
}
