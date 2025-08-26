package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

// Config holds database configuration
type Config struct {
	DBPath string
}

// DefaultConfig returns the default database configuration
func DefaultConfig() Config {
	return Config{
		DBPath: "words.db",
	}
}

// Initialize sets up the database connection
func Initialize(cfg Config) error {
	var err error
	once.Do(func() {
		// Ensure the directory exists
		dir := filepath.Dir(cfg.DBPath)
		if dir != "." {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return
			}
		}

		// Open database connection
		db, err = sql.Open("sqlite", cfg.DBPath)
		if err != nil {
			return
		}

		// Set connection parameters
		db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
		db.SetMaxIdleConns(1)

		// Test the connection
		if err = db.Ping(); err != nil {
			return
		}

		// Enable foreign key constraints
		if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
			return
		}

		log.Printf("Successfully connected to database: %s", cfg.DBPath)
	})

	return err
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	if db == nil {
		panic("Database not initialized. Call Initialize() first")
	}
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// Transaction executes a function within a database transaction
func Transaction(fn func(*sql.Tx) error) error {
	tx, err := GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
