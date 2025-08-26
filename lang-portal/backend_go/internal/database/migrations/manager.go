package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Manager handles database migrations
type Manager struct {
	db   *sql.DB
	path string
}

// NewManager creates a new migration manager
func NewManager(db *sql.DB, migrationsPath string) *Manager {
	return &Manager{
		db:   db,
		path: migrationsPath,
	}
}

// Initialize creates the migrations table if it doesn't exist
func (m *Manager) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := m.db.Exec(query)
	return err
}

// getMigrationsFiles returns a sorted list of migration files
func (m *Manager) getMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

// getAppliedMigrations returns a map of applied migrations
func (m *Manager) getAppliedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT name FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, rows.Err()
}

// Migrate runs all pending migrations
func (m *Manager) Migrate() error {
	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	files, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, file := range files {
		if applied[file] {
			continue
		}

		log.Printf("Applying migration: %s", file)
		
		content, err := ioutil.ReadFile(filepath.Join(m.path, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", file); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		log.Printf("Successfully applied migration: %s", file)
	}

	return nil
}
